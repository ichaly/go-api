package app

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

const (
	pemSig = "--BEGIN "
)

const (
	logLevelNone int = iota
	logLevelInfo
	logLevelWarn
	logLevelError
	logLevelDebug
)

type dbConf struct {
	driverName string
	connString string
}

func NewDB(conf *Config, openDB bool, log *zap.SugaredLogger, fs afero.Fs) (*sql.DB, error) {
	return newDB(conf, openDB, false, log, fs)
}

func newDB(
	conf *Config,
	openDB, useTelemetry bool,
	log *zap.SugaredLogger,
	fs afero.Fs,
) (*sql.DB, error) {

	var db *sql.DB
	var dc *dbConf
	var err error

	switch conf.DataSource.Type {
	case "mysql":
		dc, err = initMysql(conf, openDB, useTelemetry, fs)
	default:
		dc, err = initPostgres(conf, openDB, useTelemetry, fs)
	}

	if err != nil {
		return nil, fmt.Errorf("database init: %v", err)
	}

	for i := 0; ; {
		if db, err = sql.Open(dc.driverName, dc.connString); err == nil {
			db.SetMaxIdleConns(conf.DataSource.PoolSize)
			db.SetMaxOpenConns(conf.DataSource.MaxConnections)
			db.SetConnMaxIdleTime(conf.DataSource.MaxConnIdleTime)
			db.SetConnMaxLifetime(conf.DataSource.MaxConnLifeTime)

			if err := db.Ping(); err == nil {
				return db, nil
			} else {
				_ = db.Close()
				log.Warnf("database ping: %s", err)
			}

		} else {
			log.Warnf("database open: %s", err)
		}

		time.Sleep(time.Duration(i*100) * time.Millisecond)

		if i > 50 {
			return nil, err
		} else {
			i++
		}
	}
}

func initPostgres(conf *Config, openDB, useTelemetry bool, fs afero.Fs) (*dbConf, error) {
	c := conf
	config, _ := pgx.ParseConfig("")
	config.Host = c.DataSource.Host
	config.Port = c.DataSource.Port
	config.User = c.DataSource.Username
	config.Password = c.DataSource.Password

	config.RuntimeParams = map[string]string{
		"application_name": c.AppName,
		"search_path":      c.DataSource.Schema,
	}

	if openDB {
		config.Database = c.DataSource.Name
	}

	if c.DataSource.EnableTLS {
		if len(c.DataSource.ServerName) == 0 {
			return nil, errors.New("tls: server_name is required")
		}
		if len(c.DataSource.ServerCert) == 0 {
			return nil, errors.New("tls: server_cert is required")
		}

		rootCertPool := x509.NewCertPool()
		var pem []byte
		var err error

		if strings.Contains(c.DataSource.ServerCert, pemSig) {
			pem = []byte(strings.ReplaceAll(c.DataSource.ServerCert, `\n`, "\n"))
		} else {
			pem, err = afero.ReadFile(fs, c.DataSource.ServerCert)
		}

		if err != nil {
			return nil, fmt.Errorf("tls: %w", err)
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			return nil, errors.New("tls: failed to append pem")
		}

		config.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
			RootCAs:    rootCertPool,
			ServerName: c.DataSource.ServerName,
		}

		if len(c.DataSource.ClientCert) > 0 {
			if len(c.DataSource.ClientKey) == 0 {
				return nil, errors.New("tls: client_key is required")
			}

			clientCert := make([]tls.Certificate, 0, 1)
			var certs tls.Certificate

			if strings.Contains(c.DataSource.ClientCert, pemSig) {
				certs, err = tls.X509KeyPair(
					[]byte(strings.ReplaceAll(c.DataSource.ClientCert, `\n`, "\n")),
					[]byte(strings.ReplaceAll(c.DataSource.ClientKey, `\n`, "\n")),
				)
			} else {
				certs, err = loadX509KeyPair(fs, c.DataSource.ClientCert, c.DataSource.ClientKey)
			}

			if err != nil {
				return nil, fmt.Errorf("tls: %w", err)
			}

			clientCert = append(clientCert, certs)
			config.TLSConfig.Certificates = clientCert
		}
	}

	return &dbConf{"pgx", stdlib.RegisterConnConfig(config)}, nil
}

func initMysql(conf *Config, openDB, useTelemetry bool, fs afero.Fs) (*dbConf, error) {
	c := conf
	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/", c.DataSource.Username, c.DataSource.Password, c.DataSource.Host, c.DataSource.Port)

	if openDB {
		connString += c.DataSource.Name
	}

	return &dbConf{"mysql", connString}, nil
}

func loadX509KeyPair(fs afero.Fs, certFile, keyFile string) (tls.Certificate, error) {
	certPEMBlock, err := afero.ReadFile(fs, certFile)
	if err != nil {
		return tls.Certificate{}, err
	}
	keyPEMBlock, err := afero.ReadFile(fs, keyFile)
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.X509KeyPair(certPEMBlock, keyPEMBlock)
}
