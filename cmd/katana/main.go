package main

import (
	"github.com/pkg/errors"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/katana/internal/runner"
	"github.com/projectdiscovery/katana/pkg/types"
)

var (
	cfgFile string
	options = &types.Options{}
)

func main() {
	if err := readFlags(); err != nil {
		gologger.Fatal().Msgf("Could not read flags: %s\n", err)
	}

	runner, err := runner.New(options)
	if err != nil {
		gologger.Fatal().Msgf("Could not create runner: %s\n", err)
	}
	if err := runner.ExecuteCrawling(); err != nil {
		gologger.Fatal().Msgf("Could not execute crawling: %s\n", err)
	}
}

func readFlags() error {
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`Katana is a fast crawler focused on execution in automation
pipelines offering both headless and non-headless crawling.`)

	createGroup(flagSet, "input", "Input",
		flagSet.StringSliceVarP(&options.URLs, "list", "u", []string{}, "target url / list to crawl", goflags.FileCommaSeparatedStringSliceOptions),
	)

	createGroup(flagSet, "configs", "Configurations",
		flagSet.StringVar(&cfgFile, "config", "", "path to the nuclei configuration file"),
		flagSet.IntVarP(&options.MaxDepth, "depth", "d", 2, "maximum depth to crawl"),
		flagSet.IntVarP(&options.CrawlDuration, "crawl-duration", "ct", 0, "maximum duration to crawl the target for"),
		flagSet.IntVarP(&options.BodyReadSize, "max-response-size", "mrs", 2*1024*1024, "maximum response size to read"),
		flagSet.IntVar(&options.Timeout, "timeout", 10, "time to wait for request in seconds"),
		flagSet.IntVar(&options.Retries, "retries", 1, "number of times to retry the request"),
		flagSet.StringVarP(&options.Proxy, "proxy", "p", "", "http/socks5 proxy to use"),
		flagSet.RuntimeMapVarP(&options.CustomHeaders, "headers", "H", []string{}, "custom header/cookie to include in request"),
	)

	createGroup(flagSet, "filters", "Filters",
		flagSet.StringSliceVarP(&options.Scope, "crawl-scope", "cs", []string{}, "in scope target to be followed by crawler", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.StringSliceVarP(&options.OutOfScope, "crawl-out-scope", "cos", []string{}, "out of scope target to be excluded by crawler", goflags.FileCommaSeparatedStringSliceOptions),
		flagSet.BoolVarP(&options.IncludeSubdomains, "include-sub", "is", false, "include subdomains in crawl scope"),
		flagSet.StringSliceVar(&options.Extensions, "extensions", []string{}, "extensions to be explicitly allowed for crawling (* means all - default)", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&options.ExtensionsAllowList, "extensions-allow-list", []string{}, "extensions to allow from default deny list", goflags.CommaSeparatedStringSliceOptions),
		flagSet.StringSliceVar(&options.ExtensionDenyList, "extensions-deny-list", []string{}, "custom extensions for the crawl extensions deny list", goflags.CommaSeparatedStringSliceOptions),
	)

	createGroup(flagSet, "ratelimit", "Rate-Limit",
		flagSet.IntVarP(&options.Concurrency, "concurrency", "c", 300, "number of concurrent fetchers to use"),
		flagSet.IntVarP(&options.Delay, "delay", "rd", 0, "request delay between each request in seconds"),
		flagSet.IntVarP(&options.RateLimit, "rate-limit", "rl", 150, "maximum requests to send per second"),
		flagSet.IntVarP(&options.RateLimitMinute, "rate-limit-minute", "rlm", 0, "maximum number of requests to send per minute"),
	)

	createGroup(flagSet, "output", "Output",
		flagSet.StringVarP(&options.OutputFile, "output", "o", "", "file to write output to"),
		flagSet.BoolVar(&options.JSON, "json", false, "write output in JSONL(ines) format"),
		flagSet.BoolVarP(&options.NoColors, "no-color", "nc", false, "disable output content coloring (ANSI escape codes)"),
		flagSet.BoolVar(&options.Silent, "silent", false, "display output only"),
		flagSet.BoolVarP(&options.Verbose, "verbose", "v", false, "display verbose output"),
		flagSet.BoolVar(&options.Version, "version", false, "display project version"),
	)

	if err := flagSet.Parse(); err != nil {
		return errors.Wrap(err, "could not parse flags")
	}

	if cfgFile != "" {
		if err := flagSet.MergeConfigFile(cfgFile); err != nil {
			return errors.Wrap(err, "could not read config file")
		}
	}
	return nil
}

func createGroup(flagSet *goflags.FlagSet, groupName, description string, flags ...*goflags.FlagData) {
	flagSet.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}
