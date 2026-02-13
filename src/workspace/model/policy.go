package workspace_model

type Policy string

const (
	PolicyWorkspaceAdmin    Policy = "workspace.admin"
	PolicyWorkspaceSettings Policy = "workspace.settings"
	PolicyWorkspaceMembers  Policy = "workspace.members"
	PolicyPhoneConfigRead   Policy = "phone_config.read"
	PolicyPhoneConfigManage Policy = "phone_config.manage"
	PolicyContactRead       Policy = "contact.read"
	PolicyContactManage     Policy = "contact.manage"
	PolicyMessageRead       Policy = "message.read"
	PolicyMessageSend       Policy = "message.send"
	PolicyCampaignRead      Policy = "campaign.read"
	PolicyCampaignManage    Policy = "campaign.manage"
	PolicyCampaignRun       Policy = "campaign.run"
	PolicyWebhookRead       Policy = "webhook.read"
	PolicyWebhookManage     Policy = "webhook.manage"
	PolicyBillingRead       Policy = "billing.read"
	PolicyBillingManage     Policy = "billing.manage"
	PolicyBillingAdmin      Policy = "billing.admin"
)

var AllPolicies = []Policy{
	PolicyWorkspaceAdmin,
	PolicyWorkspaceSettings,
	PolicyWorkspaceMembers,
	PolicyPhoneConfigRead,
	PolicyPhoneConfigManage,
	PolicyContactRead,
	PolicyContactManage,
	PolicyMessageRead,
	PolicyMessageSend,
	PolicyCampaignRead,
	PolicyCampaignManage,
	PolicyCampaignRun,
	PolicyWebhookRead,
	PolicyWebhookManage,
	PolicyBillingRead,
	PolicyBillingManage,
	PolicyBillingAdmin,
}

func IsValidPolicy(p Policy) bool {
	for _, valid := range AllPolicies {
		if p == valid {
			return true
		}
	}
	return false
}
