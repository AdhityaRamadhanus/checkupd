checkup.config = {
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
};