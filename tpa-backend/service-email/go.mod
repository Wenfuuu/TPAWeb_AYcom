module github.com/Wenfuuu/tpa-web/service-email

go 1.24.2

replace github.com/Wenfuuu/tpa-web/pkg => ../pkg

require (
	github.com/Wenfuuu/tpa-web/pkg v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/rabbitmq/amqp091-go v1.10.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)

require gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
