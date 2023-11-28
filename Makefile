.PHONY: api 
C_INCLUDE_PATH=/home/renan/projetos/doctor_recorder/whisper.cpp/whisper.h
grpc:
	protoc api/v1/*.proto \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
	--proto_path=.
api: export C_INCLUDE_PATH=/home/renan/projetos/doctor_recorder/whisper.cpp/
api: export LIBRARY_PATH=/home/renan/projetos/doctor_recorder/whisper.cpp/
api:
	go run ./cmd/api/main.go
reload: export C_INCLUDE_PATH=/home/renan/projetos/doctor_recorder/whisper.cpp/whisper.h
reload: export LIBRARY_PATH=/home/renan/projetos/doctor_recorder/whisper.cpp/libwhisper.a
reload:
	air
trans:
	pipenv run python whisper
