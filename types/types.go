package types

type Configurations struct {
	RabbitMQConfigDetails RabbitMQConfigDetails `validate:"required" json:"RabbitMQConfigDetails"`
}

type RabbitMQConfigDetails struct {
	RabbitMQHost          string `validate:"required" json:"rabbitmq_host"`
	RabbitMQPort          int    `validate:"required" json:"rabbitmq_port"`
	RabbitMQExchange      string `validate:"required" json:"rabbitmq_exchange"`
	RabbitMQVHost         string `validate:"required" json:"rabbitmq_vhost"`
	RabbitMQDurable       bool   `validate:"required" json:"rabbitmq_durable"`
	RabbitMQUserName      string `validate:"required" json:"rabbitmq_username"`
	RabbitMQPassword      string `validate:"required" json:"rabbitmq_password"`
	RabbitMQAddTimeout    int    `validate:"required" json:"rabbitmq_add_timeout"`
	RabbitMQPrefetchCount int    `validate:"required" json:"rabbitmq_prefetch_count"`
	RabbitMQExchangeType  string `validate:"required" json:"rabbitmq_exchange_type"`

	RabbitMQComplaintResponseQueue string `validate:"required" json:"rabbitmq_complaint_response_queue"`
	RabbitMQCDRLookupQueue         string `validate:"required" json:"rabbitmq_cdr_lookup_queue"`
	RabbitMQCDRConfirmationQueue   string `validate:"required" json:"rabbitmq_cdr_confirmation_queue"`

	RabbitMQCACert                    string `validate:"required" json:"rabbitmq_cacert"`
	RabbitMQClientCert                string `validate:"required" json:"rabbitmq_clientcert"`
	RabbitMQClientKey                 string `validate:"required" json:"rabbitmq_clientkey"`
	RabbitMQSSLFlag                   int    `json:"rabbitmq_ssl_flag"`
	RABBITMQ_TLS_INSECURE_SKIP_VERIFY bool   `json:"RABBITMQ_TLS_INSECURE_SKIP_VERIFY"`
	RABBITMQ_TLS_ONLY_CA_CERT         bool   `json:"RABBITMQ_TLS_ONLY_CA_CERT"`
}