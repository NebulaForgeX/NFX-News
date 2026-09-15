#!/usr/bin/env python3
"""Patch NFX-News wirings, configs, compose, and copy otelx."""
from __future__ import annotations

import pathlib
import re
import shutil

ROOT = pathlib.Path("/volume1/Projects/NebulaForgeX/NFX-News")
IDENTITY = pathlib.Path("/volume1/Projects/NebulaForgeX/NFX-Identity")

ADAPTER_OLD = '''	return &token.Claims{Registered: claims.RegisteredClaims, Raw: map[string]any{"user_id": claims.UserID}}, nil'''
ADAPTER_NEW = '''	return &token.Claims{Registered: claims.RegisteredClaims, Raw: map[string]any{
		"account_id":    claims.AccountID,
		"profile_id":    claims.ProfileID,
		"profile_scope": claims.ProfileScope,
	}}, nil'''

OTEL_TOML = '''
[otel]
    enabled = false
    endpoint = "${OTEL_EXPORTER_OTLP_ENDPOINT}"
    insecure = true
    sampler_arg = 1.0
    traces = true
    metrics = true
    logs = true
    export_timeout = "10s"
'''

OTEL_IMPORT = '\t"nfxnews/pkgs/connections/otelx"\n'


def patch_wiring(path: pathlib.Path) -> None:
    text = path.read_text()
    if "claims.UserID" in text:
        text = text.replace(ADAPTER_OLD, ADAPTER_NEW)
        if ADAPTER_OLD in text or "claims.UserID" in text:
            text = re.sub(
                r'return &token\.Claims\{Registered: claims\.RegisteredClaims, Raw: map\[string\]any\{"user_id": claims\.UserID\}\}, nil',
                ADAPTER_NEW.strip(),
                text,
            )
    module = path.parent.parent.name  # modules/<mod>/server/wiring.go
    if "authconn.Dial" not in text:
        if 'authconn "nfxnews/connections/auth"' not in text:
            text = text.replace(
                '"nfxnews/pkgs/cachex"',
                'authconn "nfxnews/connections/auth"\n\t"nfxnews/pkgs/cachex"',
                1,
            )
        if "identityAuth" not in text:
            text = text.replace(
                "\tconns               []*grpc.ClientConn\n}",
                "\tconns               []*grpc.ClientConn\n\tidentityAuth        *authconn.Client\n}",
                1,
            )
        insert = f'''	identityClient, err := authconn.Dial(authconn.GRPCConfig{{
		Addr:           cfg.GRPCClient.AuthAddr,
		TokenSecretKey: cfg.Token.SecretKey,
		TokenIssuer:    cfg.Token.Issuer,
		CallerService:  "{module}",
	}})
	if err != nil {{
		return nil, fmt.Errorf("dial identity auth: %w", err)
	}}
'''
        if "errorsLangsPath :=" in text:
            text = text.replace("	errorsLangsPath :=", insert + "	errorsLangsPath :=", 1)
        text = text.replace(
            "		userTokenVerifier: userTokenVerifier, serverTokenVerifier: serverTokenVerifier, errorsLangsPath: errorsLangsPath,\n",
            "		userTokenVerifier: userTokenVerifier, serverTokenVerifier: serverTokenVerifier, errorsLangsPath: errorsLangsPath,\n\t\tidentityAuth: identityClient,\n",
            1,
        )
        if "d.identityAuth.Close" not in text:
            text = text.replace(
                "	for _, c := range d.conns {\n		_ = c.Close()\n	}\n}",
                "	if d.identityAuth != nil {\n		_ = d.identityAuth.Close()\n	}\n	for _, c := range d.conns {\n		_ = c.Close()\n	}\n}",
                1,
            )
        if "func (d *Dependencies) AuthClient()" not in text:
            text = text.replace(
                "func (d *Dependencies) ErrorsLangsPath() string                { return d.errorsLangsPath }",
                "func (d *Dependencies) ErrorsLangsPath() string                { return d.errorsLangsPath }\nfunc (d *Dependencies) AuthClient() *authconn.Client           { return d.identityAuth }",
                1,
            )
    path.write_text(text)


def patch_config_types(path: pathlib.Path) -> None:
    text = path.read_text()
    if "otelx.Config" in text:
        return
    if '"nfxnews/pkgs/connections/otelx"' not in text:
        text = text.replace(
            '"nfxnews/pkgs/cachex"',
            '"nfxnews/pkgs/cachex"\n\t"nfxnews/pkgs/connections/otelx"',
            1,
        )
    text = text.replace(
        "\tI18n           I18nConfig         `koanf:\"i18n\"`\n}",
        "\tI18n           I18nConfig         `koanf:\"i18n\"`\n\tOTEL           otelx.Config       `koanf:\"otel\"`\n}",
        1,
    )
    path.write_text(text)


def patch_toml(path: pathlib.Path) -> None:
    text = path.read_text()
    if "[otel]" not in text:
        path.write_text(text.rstrip() + "\n" + OTEL_TOML)


def strip_auth_compose(path: pathlib.Path) -> None:
    text = path.read_text()
    # Remove YAML service blocks that start with auth-
    text = re.sub(
        r"\n  auth-[a-z]+:\n(?:    .*\n)+?(?=\n  [a-z]|\nnetworks:)",
        "\n",
        text,
    )
    path.write_text(text)


def copy_otelx() -> None:
    src = IDENTITY / "pkgs/connections/otelx"
    dst = ROOT / "pkgs/connections/otelx"
    dst.parent.mkdir(parents=True, exist_ok=True)
    if dst.exists():
        shutil.rmtree(dst)
    shutil.copytree(src, dst)
    for p in dst.rglob("*.go"):
        p.write_text(p.read_text().replace("nfxidentity/", "nfxnews/"))


def main() -> None:
    copy_otelx()
    for wiring in (ROOT / "modules").glob("*/server/wiring.go"):
        patch_wiring(wiring)
        print("patched wiring", wiring)
    for types in (ROOT / "modules").glob("*/config/types.go"):
        patch_config_types(types)
        print("patched types", types)
    for toml in (ROOT / "inputs").glob("*/configuration/*.toml"):
        patch_toml(toml)
    strip_auth_compose(ROOT / "docker-compose.yml")
    strip_auth_compose(ROOT / "docker-compose.dev.yml")
    auth_schema = ROOT / "databases/src/schemas/auth"
    if auth_schema.exists():
        shutil.rmtree(auth_schema)
        print("removed auth schema")


if __name__ == "__main__":
    main()
