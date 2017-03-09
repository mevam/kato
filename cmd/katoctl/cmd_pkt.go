package katoctl

//-----------------------------------------------------------------------------
// 'katoctl pkt' command flags definitions:
//-----------------------------------------------------------------------------

var (

	//------------------------
	// pkt: top level command
	//------------------------

	cmdPkt = app.Command("pkt", "Kato's Packet.net provider.")

	//----------------------------
	// pkt deploy: nested command
	//----------------------------

	cmdPktDeploy = cmdPkt.Command("deploy",
		"Deploy Kato's infrastructure on Packet.net")

	//---------------------------
	// pkt setup: nested command
	//---------------------------

	cmdPktSetup = cmdPkt.Command("setup",
		"Setup a Packet.net project to be used by katoctl.")

	//-------------------------
	// pkt run: nested command
	//-------------------------

	cmdPktRun = cmdPkt.Command("run", "Starts a CoreOS instance on Packet.net.")

	flPktRunAPIKey = cmdPktRun.Flag("api-key",
		"Packet API key.").
		Required().PlaceHolder("KATO_PKT_RUN_APIKEY").
		OverrideDefaultFromEnvar("KATO_PKT_RUN_APIKEY").
		String()

	flPktRunHostName = cmdPktRun.Flag("host-name",
		"Used in the Packet.net dashboard.").
		Required().PlaceHolder("KATO_PKT_RUN_HOST_NAME").
		OverrideDefaultFromEnvar("KATO_PKT_RUN_HOST_NAME").
		String()

	flPktRunProjectID = cmdPktRun.Flag("project-id",
		"Format: aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee").
		Required().PlaceHolder("KATO_PKT_RUN_PROJECT_ID").
		OverrideDefaultFromEnvar("KATO_PKT_RUN_PROJECT_ID").
		String()

	flPktRunPlan = cmdPktRun.Flag("plan",
		"One of [ baremetal_0 | baremetal_1 | baremetal_2 | baremetal_3 ]").
		Required().PlaceHolder("KATO_PKT_RUN_PLAN").
		OverrideDefaultFromEnvar("KATO_PKT_RUN_PLAN").
		Enum("baremetal_0", "baremetal_1", "baremetal_2", "baremetal_3")

	flPktRunOS = cmdPktRun.Flag("os",
		"One of [ coreos_stable | coreos_beta | coreos_alpha ]").
		Default("coreos_stable").OverrideDefaultFromEnvar("KATO_PKT_RUN_OS").
		Enum("coreos_stable", "coreos_beta", "coreos_alpha")

	flPktRunFacility = cmdPktRun.Flag("facility",
		"One of [ ewr1 | ams1 ]").
		Required().PlaceHolder("KATO_PKT_RUN_FACILITY").
		OverrideDefaultFromEnvar("KATO_PKT_RUN_FACILITY").
		Enum("ewr1", "ams1")

	flPktRunBilling = cmdPktRun.Flag("billing",
		"One of [ hourly | monthly ]").
		Default("hourly").OverrideDefaultFromEnvar("KATO_PKT_RUN_BILLING").
		Enum("hourly", "monthly")
)
