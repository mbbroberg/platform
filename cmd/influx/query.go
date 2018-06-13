package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/influxdata/platform"
	"github.com/influxdata/platform/http"
	"github.com/influxdata/platform/query"
	_ "github.com/influxdata/platform/query/builtin"
	"github.com/influxdata/platform/query/execute"
	"github.com/influxdata/platform/query/functions"
	"github.com/influxdata/platform/query/functions/storage"
	"github.com/influxdata/platform/query/functions/storage/pb"
	qid "github.com/influxdata/platform/query/id"
	"github.com/influxdata/platform/query/repl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var queryCmd = &cobra.Command{
	Use:   "query [query literal or @/path/to/query.flux]",
	Short: "Execute an Flux query",
	Long: `Execute a literal Flux query provided as a string,
		or execute a literal Flux query contained in a file by specifying the file prefixed with an @ sign.`,
	Args: cobra.ExactArgs(1),
	Run:  fluxQueryF,
}

var queryFlags struct {
	StorageHosts string
	OrgID        string
	Verbose      bool
}

func init() {
	queryCmd.PersistentFlags().StringVar(&queryFlags.StorageHosts, "storage-hosts", "localhost:8082", "Comma-separated list of storage hosts")
	viper.BindEnv("STORAGE_HOSTS")
	if h := viper.GetString("STORAGE_HOSTS"); h != "" {
		queryFlags.StorageHosts = h
	}

	queryCmd.PersistentFlags().BoolVarP(&queryFlags.Verbose, "verbose", "v", false, "Verbose output")
	viper.BindEnv("VERBOSE")
	if viper.GetBool("VERBOSE") {
		queryFlags.Verbose = true
	}

	queryCmd.PersistentFlags().StringVar(&queryFlags.OrgID, "org-id", "", "Organization ID")
	viper.BindEnv("ORG_ID")
	if h := viper.GetString("ORG_ID"); h != "" {
		queryFlags.OrgID = h
	}
}

func fluxQueryF(cmd *cobra.Command, args []string) {
	q, err := repl.LoadQuery(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	hosts, err := storageHostReader(strings.Split(queryFlags.StorageHosts, ","))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	org, err := orgID(queryFlags.OrgID)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	buckets, err := bucketService(flags.host, flags.token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	orgs, err := orgService(flags.host, flags.token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	r, err := getFluxREPL(hosts, buckets, orgs, org, queryFlags.Verbose)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := r.Input(q); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func injectDeps(deps execute.Dependencies, hosts storage.Reader, buckets platform.BucketService, orgs platform.OrganizationService) error {
	return functions.InjectFromDependencies(deps, storage.Dependencies{
		Reader:             hosts,
		BucketLookup:       query.FromBucketService(buckets),
		OrganizationLookup: query.FromOrganizationService(orgs),
	})
}

func storageHostReader(hosts []string) (storage.Reader, error) {
	return pb.NewReader(storage.NewStaticLookup(hosts))
}

func bucketService(addr, token string) (platform.BucketService, error) {
	if addr == "" {
		return nil, fmt.Errorf("bucket host address required")
	}

	return &http.BucketService{
		Addr:  addr,
		Token: token,
	}, nil
}

func orgService(addr, token string) (platform.OrganizationService, error) {
	if addr == "" {
		return nil, fmt.Errorf("organization host address required")
	}

	return &http.OrganizationService{
		Addr:  addr,
		Token: token,
	}, nil
}

func orgID(org string) (qid.ID, error) {
	var oid qid.ID
	err := oid.DecodeFromString(org)
	return oid, err
}
