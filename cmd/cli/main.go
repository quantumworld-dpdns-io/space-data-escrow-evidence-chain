package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/quantumworld-dpdns-io/space-data-escrow-evidence-chain/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}
	baseURL := envOr("EVIDENCE_API_URL", "http://localhost:8080")
	apiKey := envOr("EVIDENCE_API_KEY", "dev-api-key")
	c := cli.New(baseURL, apiKey)

	cmd := os.Args[1]
	switch cmd {
	case "evidence-create":
		if len(os.Args) < 6 {
			fmt.Println("usage: evidence-create <external_id> <source> <type> <k=v,k2=v2>")
			os.Exit(1)
		}
		payload := parseKV(os.Args[5])
		out, _, err := c.DoEvidenceCreate(os.Args[2], os.Args[3], os.Args[4], payload)
		exitOnErr(err)
		printJSON(out)
	case "evidence-verify":
		if len(os.Args) < 3 {
			fmt.Println("usage: evidence-verify <evidence_id>")
			os.Exit(1)
		}
		out, _, err := c.DoEvidenceVerify(os.Args[2])
		exitOnErr(err)
		printJSON(out)
	case "audit-query":
		out, _, err := c.DoAuditQuery()
		exitOnErr(err)
		printJSON(out)
	case "enrich-trigger":
		if len(os.Args) < 3 {
			fmt.Println("usage: enrich-trigger <evidence_id>")
			os.Exit(1)
		}
		out, _, err := c.DoEnrichTrigger(os.Args[2])
		exitOnErr(err)
		printJSON(out)
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Println("commands: evidence-create, evidence-verify, audit-query, enrich-trigger")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseKV(raw string) map[string]string {
	out := map[string]string{}
	if raw == "" {
		return out
	}
	for _, p := range strings.Split(raw, ",") {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 2 {
			out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return out
}

func printJSON(b []byte) {
	var obj any
	if err := json.Unmarshal(b, &obj); err != nil {
		fmt.Println(string(b))
		return
	}
	pretty, _ := json.MarshalIndent(obj, "", "  ")
	fmt.Println(string(pretty))
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
