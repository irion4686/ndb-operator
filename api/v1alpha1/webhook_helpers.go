package v1alpha1

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/nutanix-cloud-native/ndb-operator/api"
	"github.com/nutanix-cloud-native/ndb-operator/common"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Get specific implementation of the DBProvisionRequestAppender interface based on the provided databaseType
func getDatabaseWebhookHandler(database *Database) DatabaseWebhookHandler {
	if database.Spec.IsClone {
		return &CloningHandler{}
	} else {
		return &ProvisoningHandler{}
	}
}

// +kubebuilder:object:generate:=false
type DatabaseWebhookHandler interface {
	// Default logic
	defaulter(databaseSpec *DatabaseSpec)
	// Validates creation (after defaulting)
	validateCreate(databaseSpec *DatabaseSpec, errors *field.ErrorList, instancePath *field.Path)
}

// +kubebuilder:object:generate:=false
// Implements Validator
type CloningHandler struct{}

// +kubebuilder:object:generate:=false
// Implements Validator
type ProvisoningHandler struct{}

func (v *CloningHandler) defaulter(spec *DatabaseSpec) {
	databaselog.Info("Entering Clone defaulter logic...")

	// Default Database
	if spec.Instance == nil {
		spec.Instance = &(Instance{})
	}
	spec.Instance.DatabaseNames = []string{}
	spec.Instance.Profiles = &(Profiles{})
	spec.Instance.TMInfo = &(DBTimeMachineInfo{})
	spec.Instance.AdditionalArguments = map[string]string{}

	// Default Clone
	if spec.Clone == nil {
		spec.Clone = &(Clone{})
	}

	if spec.Clone.Description == "" {
		description := "Clone created by ndb-operator: " + spec.Clone.Name
		databaselog.Info(fmt.Sprintf("Initializing Description to: %s.", description))
		spec.Clone.Description = description
	}

	if spec.Clone.Profiles == nil {
		databaselog.Info(fmt.Sprintf("Initializing empty Profiles"))
		spec.Clone.Profiles = &(Profiles{})
	}

	if spec.Clone.AdditionalArguments == nil {
		databaselog.Info(fmt.Sprintf("Initializing Description to empty map."))
		spec.Clone.AdditionalArguments = map[string]string{}
	}

	databaselog.Info("Exiting Clone defaulter logic!")
}

func (v *CloningHandler) validateCreate(spec *DatabaseSpec, errors *field.ErrorList, clonePath *field.Path) {
	databaselog.Info("Entering Clone defaulter logic...")
	databaselog.Info("Exiting Clone defaulter logic!")

	clone := spec.Clone

	if clone.Name == "" {
		*errors = append(*errors, field.Invalid(clonePath.Child("name"), clone.Name, "A valid Name must be specified"))
	}

	if clone.ClusterId == "" {
		*errors = append(*errors, field.Invalid(clonePath.Child("clusterId"), clone.ClusterId, "Ensure ClusterId is a valid UUID"))
	}

	if clone.CredentialSecret == "" {
		*errors = append(*errors, field.Invalid(clonePath.Child("credentialSecret"), clone.CredentialSecret, "CredentialSecret must be provided in the Clone Spec"))
	}

	if clone.TimeZone == "" {
		*errors = append(*errors, field.Invalid(clonePath.Child("timeZone"), clone.CredentialSecret, "CredentialSecret must be provided in Clone Spec"))
	}

	if clone.SourceDatabaseId == "" {
		*errors = append(*errors, field.Invalid(clonePath.Child("sourceDatabaseId"), clone.CredentialSecret, "sourceDatabaseId must be provided"))
	}

	if clone.SnapshotId == "" {
		*errors = append(*errors, field.Invalid(clonePath.Child("sourceDatabaseId"), clone.CredentialSecret, "snapshotId must be provided"))
	}

	if _, isPresent := api.AllowedDatabaseTypes[clone.Type]; !isPresent {
		*errors = append(*errors, field.Invalid(clonePath.Child("type"), clone.Type,
			fmt.Sprintf("A valid clone type must be specified. Valid values are: %s", reflect.ValueOf(api.AllowedDatabaseTypes).MapKeys()),
		))
	}

	if _, isPresent := api.ClosedSourceDatabaseTypes[clone.Type]; isPresent {
		if clone.Profiles == &(Profiles{}) || clone.Profiles.Software == (Profile{}) {
			*errors = append(*errors, field.Invalid(clonePath.Child("profiles").Child("software"), clone.Profiles.Software, "Software Profile must be provided for the closed-source database engines"))
		}
	}

	if err := additionalArgumentsValidationCheck(spec.IsClone, clone.Type, clone.AdditionalArguments); err != nil {
		*errors = append(*errors, field.Invalid(clonePath.Child("additionalArguments"), clone.AdditionalArguments, err.Error()))
	}

}

func (v *ProvisoningHandler) defaulter(spec *DatabaseSpec) {
	databaselog.Info("Entering Database defaulter logic...")

	// Default Clone
	if spec.Clone == nil {
		spec.Clone = &(Clone{})
	}
	spec.Clone.Profiles = &(Profiles{})
	spec.Clone.AdditionalArguments = map[string]string{}

	// Default Database
	if spec.Instance == nil {
		spec.Instance = &(Instance{})
	}

	if spec.Instance.Description == "" {
		description := "Database provisioned by ndb-operator: " + spec.Instance.Name
		databaselog.Info(fmt.Sprintf("Initializing Description to: %s.", description))
		spec.Instance.Description = description
	}

	if len(spec.Instance.DatabaseNames) == 0 {
		databaselog.Info(fmt.Sprintf("Initializing DatabaseNames to: %s.", api.DefaultDatabaseNames))
		spec.Instance.DatabaseNames = api.DefaultDatabaseNames
	}

	if spec.Instance.TimeZone == "" {
		databaselog.Info(fmt.Sprintf("Initializing TimeZone to: %s.", common.TIMEZONE_UTC))
		spec.Instance.TimeZone = common.TIMEZONE_UTC
	}

	// Initialize Profiles field, if it's nil (mandatory)
	if spec.Instance.Profiles == nil {
		databaselog.Info("Initializing empty Profiles.")
		spec.Instance.Profiles = &(Profiles{})
	}

	// Profiles defaulting logic
	if spec.Instance.Profiles.Compute.Id == "" && spec.Instance.Profiles.Compute.Name == "" {
		databaselog.Info(fmt.Sprintf("Initializing Profiles.Compute.Name to: %s", common.PROFILE_DEFAULT_OOB_SMALL_COMPUTE))
		spec.Instance.Profiles.Compute.Name = common.PROFILE_DEFAULT_OOB_SMALL_COMPUTE
	}

	// Initialize TMInfo field, if it's nil (mandatory)
	if spec.Instance.TMInfo == nil {
		databaselog.Info("Initializing empty tmInfo.")
		spec.Instance.TMInfo = &(DBTimeMachineInfo{})
	}

	// TMInfo defaulting logic
	if spec.Instance.TMInfo.SnapshotsPerDay == 0 {
		databaselog.Info(fmt.Sprintf("Initializing TMInfo.SnapshotsPerDay to: %d", 1))
		spec.Instance.TMInfo.SnapshotsPerDay = 1
	}

	if spec.Instance.TMInfo.SLAName == "" {
		databaselog.Info(fmt.Sprintf("Initializing TMInfo.SLAName to: %s", common.SLA_NAME_NONE))
		spec.Instance.TMInfo.SLAName = common.SLA_NAME_NONE
	}

	if spec.Instance.TMInfo.DailySnapshotTime == "" {
		databaselog.Info(fmt.Sprintf("Initializing TMInfo.DailySnapshotTime to: %s", "04:00:00"))
		spec.Instance.TMInfo.DailySnapshotTime = "04:00:00"
	}

	if spec.Instance.TMInfo.LogCatchUpFrequency == 0 {
		databaselog.Info(fmt.Sprintf("Initializing TMInfo.LogCatchUpFrequency to: %d", 30))
		spec.Instance.TMInfo.LogCatchUpFrequency = 30
	}

	if spec.Instance.TMInfo.WeeklySnapshotDay == "" {
		databaselog.Info(fmt.Sprintf("Initializing TMInfo.WeeklySnapshotDay to: %s", "FRIDAY"))
		spec.Instance.TMInfo.WeeklySnapshotDay = "FRIDAY"
	}

	if spec.Instance.TMInfo.MonthlySnapshotDay == 0 {
		databaselog.Info(fmt.Sprintf("Initializing TMInfo.MonthlySnapshotDay to: %d", 15))
		spec.Instance.TMInfo.MonthlySnapshotDay = 15
	}

	if spec.Instance.TMInfo.QuarterlySnapshotMonth == "" {
		databaselog.Info(fmt.Sprintf("Initializing TMInfo.QuarterlySnapshotMonth to: %s", "Jan"))
		spec.Instance.TMInfo.QuarterlySnapshotMonth = "Jan"
	}

	// Initialize AdditionalArguments field, if it's nil (mandatory)
	if spec.Instance.AdditionalArguments == nil {
		databaselog.Info("Initializing empty AdditionalArguments...")
		spec.Instance.AdditionalArguments = map[string]string{}
	}

	databaselog.Info("Exiting Database defaulter logic!")
}

func (v *ProvisoningHandler) validateCreate(spec *DatabaseSpec, errors *field.ErrorList, instancePath *field.Path) {
	databaselog.Info("Entering Database validateCreate logic...")

	instance := spec.Instance
	tmInfo := instance.TMInfo
	tmPath := instancePath.Child("timeMachine")

	if instance.Name == "" {
		*errors = append(*errors, field.Invalid(instancePath.Child("name"), instance.Name, "A valid Name must be specified"))
	}

	if instance.ClusterId == "" {
		*errors = append(*errors, field.Invalid(instancePath.Child("clusterId"), instance.ClusterId, "Ensure ClusterId is a valid UUID"))
	}

	if instance.Size < 10 {
		*errors = append(*errors, field.Invalid(instancePath.Child("size"), instance.Size, "Initial Database size must be specified with a value 10 GBs or more"))
	}

	if instance.CredentialSecret == "" {
		*errors = append(*errors, field.Invalid(instancePath.Child("credentialSecret"), instance.CredentialSecret, "CredentialSecret must be provided in the Instance Spec"))
	}

	if _, isPresent := api.AllowedDatabaseTypes[instance.Type]; !isPresent {
		*errors = append(*errors, field.Invalid(instancePath.Child("type"), instance.Type,
			fmt.Sprintf("A valid database type must be specified. Valid values are: %s", reflect.ValueOf(api.AllowedDatabaseTypes).MapKeys()),
		))
	}

	if _, isPresent := api.ClosedSourceDatabaseTypes[instance.Type]; isPresent {
		if instance.Profiles == &(Profiles{}) || instance.Profiles.Software == (Profile{}) {
			*errors = append(*errors, field.Invalid(instancePath.Child("profiles").Child("software"), instance.Profiles.Software, "Software Profile must be provided for the closed-source database engines"))
		}
	}

	// validating time machine info
	dailySnapshotTimeRegex := regexp.MustCompile(`^(2[0-3]|[01][0-9]):[0-5][0-9]:[0-5][0-9]$`)
	if isMatch := dailySnapshotTimeRegex.MatchString(tmInfo.DailySnapshotTime); !isMatch {
		*errors = append(*errors, field.Invalid(tmPath.Child("dailySnapshotTime"), tmInfo.DailySnapshotTime, "Invalid time format for the daily snapshot time. Use the 24-hour format (HH:MM:SS)."))
	}

	if tmInfo.SnapshotsPerDay < 1 || tmInfo.SnapshotsPerDay > 6 {
		*errors = append(*errors, field.Invalid(tmPath.Child("snapshotsPerDay"), tmInfo.SnapshotsPerDay, "Number of snapshots per day should be within 1 to 6"))
	}

	if _, isPresent := api.AllowedLogCatchupFrequencyInMinutes[tmInfo.LogCatchUpFrequency]; !isPresent {
		*errors = append(*errors, field.Invalid(tmPath.Child("logCatchUpFrequency"), tmInfo.LogCatchUpFrequency,
			fmt.Sprintf("Log catchup frequency must be specified. Valid values are: %s", reflect.ValueOf(api.AllowedLogCatchupFrequencyInMinutes).MapKeys()),
		))
	}

	if _, isPresent := api.AllowedWeeklySnapshotDays[tmInfo.WeeklySnapshotDay]; !isPresent {
		*errors = append(*errors, field.Invalid(tmPath.Child("weeklySnapshotDay"), tmInfo.WeeklySnapshotDay,
			fmt.Sprintf("Weekly Snapshot day must be specified. Valid values are: %s", reflect.ValueOf(api.AllowedWeeklySnapshotDays).MapKeys()),
		))
	}

	if tmInfo.MonthlySnapshotDay < 1 || tmInfo.MonthlySnapshotDay > 28 {
		*errors = append(*errors, field.Invalid(tmPath.Child("monthlySnapshotDay"), tmInfo.MonthlySnapshotDay, "Monthly snapshot day value must be between 1 and 28"))
	}

	if _, isPresent := api.AllowedQuarterlySnapshotMonths[tmInfo.QuarterlySnapshotMonth]; !isPresent {
		*errors = append(*errors, field.Invalid(tmPath.Child("quarterlySnapshotMonth"), tmInfo.QuarterlySnapshotMonth,
			fmt.Sprintf("Quarterly snapshot month must be specified. Valid values are: %s", reflect.ValueOf(api.AllowedQuarterlySnapshotMonths).MapKeys()),
		))
	}

	if err := additionalArgumentsValidationCheck(spec.IsClone, instance.Type, instance.AdditionalArguments); err != nil {
		*errors = append(*errors, field.Invalid(instancePath.Child("additionalArguments"), instance.AdditionalArguments, err.Error()))
	}

	databaselog.Info("Existing Database validateCreate logic!")
}
