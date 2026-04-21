package env

type Environment string

const (
	Development Environment = "development"
	Dev         Environment = "dev"
	Develop     Environment = "develop"
	Production  Environment = "production"
	Prod        Environment = "prod"
	QA          Environment = "qa"
	UAT         Environment = "uat"
)

func (e Environment) IsDevelopment() bool {
	return e == Dev || e == Develop || e == Development
}

func (e Environment) IsProduction() bool {
	return e == Prod || e == Production
}

func (e Environment) IsQA() bool {
	return e == QA || e == UAT
}

func (e Environment) String() string {
	return string(e)
}

func Parse(env string) Environment {
	switch env {
	case "dev", "develop", "development":
		return Development
	case "prod", "production":
		return Production
	case "qa", "uat":
		return QA
	default:
		return Development
	}
}
