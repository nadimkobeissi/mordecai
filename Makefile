# SPDX-FileCopyrightText: © 2019-2020 Nadim Kobeissi <nadim@symbolic.software>
# SPDX-License-Identifier: GPL-3.0-only

all:
	@go build -trimpath -gcflags="-e" -ldflags="-s -w" -o mordecai mordecai.go
	@upx -q9 mordecai
