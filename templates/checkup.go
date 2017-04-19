package templates

var ConfigFSJS = `checkup.config = {
	"timeframe": 1 * time.Day,
	"refresh_interval": 60,
	"storage": {
		"url": "logs"
	},
	"status_text": {
		"healthy": "Situation Normal",
		"degraded": "Degraded Service",
		"down": "Service Disruption"
	}
};`

var ConfigS3JS = `checkup.config = {
	"timeframe": 1 * time.Day,
	"refresh_interval": 60,
	"storage": {
		"AccessKeyID": "{{.AccessKeyID}}",
		"SecretAccessKey": "{{.SecretAccessKey}}",
		"Region": "{{.Region}}",
		"BucketName": "{{.Bucket}}"
	},
	"status_text": {
		"healthy": "Situation Normal",
		"degraded": "Degraded Service",
		"down": "Service Disruption"
	}
};`

var IndexHTML = `<!DOCTYPE html>
<html>
	<head>
		<title>Status Page</title>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<script src="js/d3.v3.min.js" charset="utf-8"></script>
		<script src="js/{{.Type}}.js"></script>
		<script src="js/checkup.js"></script>
		<script src="js/config.js"></script>
		<script src="js/statuspage.js"></script>
		<link rel="icon" href="images/favicon.png" id="favicon">
		<link rel="stylesheet" href="css/fontello.css">
		<link rel="stylesheet" href="css/style.css">
	</head>
	<body>
		<div class="app">
			<header>
				<div id="overall-status">
					<i class="icon-stethoscope overall-status-icon"></i>
					<span class="overall-status-text">Situation Normal</span>
				</div>
			</header>

			<main>
				<div class="endpoint-status">
				<div class="infobar-item border-top-black width-25">
					<span class="totalcheck-text">
						<b>Total Checks</b> <br>
						<span id="info-totalchecks">0</span>
					</span>
				</div>
				<div class="infobar-item border-top-green width-25">
					<span class="totalhealthy-text">
						<b>Healthy Endpoints</b> <br>
						<span id="info-totalhealthy">0</span>
					</span>
				</div>
				<div class="infobar-item border-top-red width-25">
					<span class="totaldown-text">
						<b>Down Endpoints</b> <br>
						<span id="info-totaldown">0</span>
					</span>
				</div>
				<div class="infobar-item border-top-black width-25">
					<span class="overall-lastcheck-text">
						<b>Last check</b> <br><span id="info-lastcheck"><time class="dynamic" datetime="">Unknown</time>
						</span>
					</span>
				</div>
				<div id="chart-grid">
					<span id="chart-placeholder">&nbsp;</span>
				</div>
			</div>
			<div id="timeline">
				<div id="big-gap">
					There is a big gap of time where no checkups were performed, so some graphs may look distorted.
				</div>
				<div id="bg-line"></div>
			</div>
		</main>

		<footer>
			Powered by <img src="images/checkup.png" id="checkup-logo">
		</footer>
	</body>
</html>`
