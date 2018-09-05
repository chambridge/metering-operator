package e2e

import (
	"flag"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/operator-framework/operator-metering/test/framework"
)

var (
	testFramework *framework.Framework

	reportTestTimeout         = 5 * time.Minute
	reportTestOutputDirectory string
	runAWSBillingTests        bool

	periodStart, periodEnd time.Time
)

func init() {
	if reportTestTimeoutStr := os.Getenv("REPORT_TEST_TIMEOUT"); reportTestTimeoutStr != "" {
		var err error
		reportTestTimeout, err = time.ParseDuration(reportTestTimeoutStr)
		if err != nil {
			log.Fatalf("Invalid REPORT_TEST_TIMEOUT: %v", err)
		}
	}
	reportTestOutputDirectory = os.Getenv("TEST_RESULT_REPORT_OUTPUT_DIRECTORY")
	if reportTestOutputDirectory == "" {
		log.Fatalf("$TEST_RESULT_REPORT_OUTPUT_DIRECTORY must be set")
	}

	err := os.MkdirAll(reportTestOutputDirectory, 0777)
	if err != nil {
		log.Fatalf("error making directory %s, err: %s", reportTestOutputDirectory, err)
	}

	runAWSBillingTests = os.Getenv("ENABLE_AWS_BILLING_TESTS") == "true"
}

func TestMain(m *testing.M) {
	kubeconfig := flag.String("kubeconfig", "", "kube config path, e.g. $HOME/.kube/config")
	ns := flag.String("namespace", "metering-ci", "test namespace")
	httpsAPI := flag.Bool("https-api", false, "If true, use https to talk to Metering API")
	flag.Parse()

	var err error
	if testFramework, err = framework.New(*ns, *kubeconfig, *httpsAPI); err != nil {
		logrus.Fatalf("failed to setup framework: %v\n", err)
	}

	os.Exit(m.Run())
}

func TestReportingE2E(t *testing.T) {
	t.Run("TestReportingProducesResults", func(t *testing.T) {
		// validate all the ReportDataSources for our tests exist before running
		// collect
		for _, test := range reportsProduceDataTestCases {
			if test.skip {
				continue
			}
			reportGenQuery, err := testFramework.WaitForMeteringReportGenerationQuery(t, test.queryName, time.Second*5, test.timeout)
			require.NoError(t, err, "ReportGenerationQuery should exist before creating report using it")

			for _, datasourceName := range reportGenQuery.Spec.DataSources {
				_, err := testFramework.WaitForMeteringReportDataSourceTable(t, datasourceName, time.Second*5, test.timeout)
				require.NoError(t, err, "ReportDataSource %s table for ReportGenerationQuery %s should exist before running reports against it", datasourceName, test.queryName)
			}
		}

		for _, test := range scheduledReportsProduceDataTestCases {
			reportGenQuery, err := testFramework.WaitForMeteringReportGenerationQuery(t, test.queryName, time.Second*5, test.timeout)
			require.NoError(t, err, "ReportGenerationQuery should exist before creating report using it")

			for _, datasourceName := range reportGenQuery.Spec.DataSources {
				_, err := testFramework.WaitForMeteringReportDataSourceTable(t, datasourceName, time.Second*5, test.timeout)
				require.NoError(t, err, "ReportDataSource %s table for ReportGenerationQuery %s should exist before running reports against it", datasourceName, test.queryName)
			}
		}

		periodStart, periodEnd = testFramework.CollectMetricsOnce(t)

		t.Run("TestReportsProduceData", testReportsProduceData)
		t.Run("TestScheduledReportsProduceData", testScheduledReportsProduceData)
	})
}
