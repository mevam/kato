package udata

//-----------------------------------------------------------------------------
// Package factored import statement:
//-----------------------------------------------------------------------------

import (

	// Stdlib:
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"

	// Community:
	log "github.com/Sirupsen/logrus"
)

//-----------------------------------------------------------------------------
// Typedefs:
//-----------------------------------------------------------------------------

// Data contains variables to be interpolated in templates.
type Data struct {
	HostID           string
	Domain           string
	Role             string
	Ns1ApiKey        string
	CaCert           string
	EtcdToken        string
	GzipUdata        bool
	FlannelNetwork   string
	FlannelSubnetLen string
	FlannelSubnetMin string
	FlannelSubnetMax string
	FlannelBackend   string
}

//-----------------------------------------------------------------------------
// func: etcdToken
//-----------------------------------------------------------------------------

func (d *Data) etcdToken() error {

	if d.EtcdToken == "auto" {

		// Request an etcd bootstrap token:
		res, err := http.Get("https://discovery.etcd.io/new?size=3")
		if err != nil {
			log.WithField("cmd", "udata").Error(err)
			return err
		}

		// Retrieve the token URL:
		tokenURL, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.WithField("cmd", "udata").Error(err)
			return err
		}

		// Extract the token ID:
		slice := strings.Split(string(tokenURL), "/")
		d.EtcdToken = slice[len(slice)-1]
	}

	return nil
}

//-----------------------------------------------------------------------------
// func: caCert
//-----------------------------------------------------------------------------

func (d *Data) caCert() error {

	if d.CaCert != "" {

		data, err := ioutil.ReadFile(d.CaCert)
		if err != nil {
			log.WithField("cmd", "udata").Error(err)
			return err
		}

		d.CaCert = strings.TrimSpace(strings.
			Replace(string(data), "\n", "\n    ", -1))
	}

	return nil
}

//-----------------------------------------------------------------------------
// func: Render
//-----------------------------------------------------------------------------

// Render takes a Data structure and outputs valid CoreOS cloud-config
// in YAML format to stdout.
func (d *Data) Render() error {

	var err error

	// Read the CA certificate:
	if err = d.caCert(); err != nil {
		return err
	}

	// Retrieve the etcd token:
	if err = d.etcdToken(); err != nil {
		return err
	}

	// Role-based parsing:
	t := template.New("udata")

	switch d.Role {
	case "master":
		t, err = t.Parse(templMaster)
	case "node":
		t, err = t.Parse(templNode)
	case "edge":
		t, err = t.Parse(templEdge)
	}

	if err != nil {
		log.WithField("cmd", "udata").Error(err)
		return err
	}

	// Apply parsed template to data object:
	if d.GzipUdata {
		log.WithField("cmd", "udata").Info("- Rendering gzipped cloud-config template")
		w := gzip.NewWriter(os.Stdout)
		defer w.Close()
		if err = t.Execute(w, d); err != nil {
			log.WithField("cmd", "udata").Error(err)
			return err
		}
	} else {
		log.WithField("cmd", "udata").Info("- Rendering plain text cloud-config template")
		if err = t.Execute(os.Stdout, d); err != nil {
			log.WithField("cmd", "udata").Error(err)
			return err
		}
	}

	// Return on success:
	return nil
}
