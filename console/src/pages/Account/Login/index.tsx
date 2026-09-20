import { AnimatedIcon, ArrowNarrowRightIcon, EyeIcon, EyeOffIcon } from "nfx-ui/icons";
import type { Login } from "nfx-ui/types";

import { useMemo, useRef, useState } from "react";
import { useGSAP } from "@gsap/react";
import { Button, Checkbox, Flex, Link, Text, TextField } from "@radix-ui/themes";
import gsap from "gsap";
import { APP_NAME } from "nfx-ui/config";
import { useLoginWithEmail } from "nfx-ui/hooks";
import { LoginFormData, useInitLoginForm } from "nfx-ui/schemas";
import { Controller, FormProvider, SubmitHandler } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { routerEventEmitter } from "@/events/router";
import { ROUTES } from "@/navigations";
import AuthToolbar from "@/pages/Account/shared/AuthToolbar";
import { safeArray, safeOr } from "@/utils";

import SelectDesk from "./SelectDesk";
import styles from "./s.module.css";

gsap.registerPlugin(useGSAP);

type WireItem = { title: string; meta: string };

export default function LoginPage() {
  const { t, i18n } = useTranslation("pages.Account.Login");
  const form = useInitLoginForm();
  const login = useLoginWithEmail();
  const [profiles, setProfiles] = useState<Login.ProfileItem[]>([]);
  const [showPassword, setShowPassword] = useState(false);
  const pageRef = useRef<HTMLDivElement>(null);

  useGSAP(
    () => {
      if (profiles.length > 0) return;
      if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
      gsap.set(".js-col", { autoAlpha: 0, x: -32 });
      gsap.set(".js-field", { autoAlpha: 0, y: 10 });
      const tl = gsap.timeline({ defaults: { ease: "power3.out" } });
      tl.to(".js-col", { autoAlpha: 1, x: 0, duration: 0.5, stagger: 0.08 }).to(".js-field", { autoAlpha: 1, y: 0, duration: 0.4, stagger: 0.06 }, "-=0.28");
    },
    { scope: pageRef, dependencies: [profiles.length] },
  );

  const onSubmit: SubmitHandler<LoginFormData> = async (data) => {
    const result = await login.mutateAsync({
      email: data.email,
      password: data.password,
      rememberMe: safeOr(data.rememberMe, false),
    });
    const list = safeArray(result?.profiles);
    if (list.length > 0) {
      setProfiles(list);
      return;
    }
    routerEventEmitter.navigate({ to: ROUTES.USER_OVERVIEW, replace: true });
  };

  const dateLabel = useMemo(
    () =>
      new Intl.DateTimeFormat(i18n.language, {
        weekday: "short",
        day: "numeric",
        month: "short",
        year: "numeric",
      }).format(new Date()),
    [i18n.language],
  );

  const worldRaw = t("wire.world.items", { returnObjects: true });
  const deskRaw = t("wire.desk.items", { returnObjects: true });
  const worldItems = Array.isArray(worldRaw) ? (worldRaw as WireItem[]) : [];
  const deskItems = Array.isArray(deskRaw) ? (deskRaw as WireItem[]) : [];

  if (profiles.length > 0) {
    return <SelectDesk profiles={profiles} onBack={() => setProfiles([])} />;
  }

  return (
    <main ref={pageRef} className={styles.page}>
      <header className={styles.masthead}>
        <div className={styles.brand}>
          <span className={styles.wordmark}>{APP_NAME}</span>
          <time className={styles.dateline} dateTime={new Date().toISOString()}>
            {dateLabel}
          </time>
          <span className={styles.edition}>{t("masthead.edition")}</span>
        </div>
        <AuthToolbar />
      </header>

      <div className={styles.board}>
        <section className={`${styles.column} ${styles.wire} js-col`}>
          <div className={styles.columnHeader}>
            <span className={styles.rail} style={{ background: "var(--blue-9)" }} />
            <div className={styles.headText}>
              <span className={styles.headName}>{t("wire.world.name")}</span>
              <span className={styles.headMeta}>{t("wire.world.meta")}</span>
            </div>
          </div>
          <div className={styles.list}>
            {worldItems.map((item, index) => (
              <div key={item.title} className={styles.item}>
                <span className={styles.rank}>{String(index + 1).padStart(2, "0")}</span>
                <span className={styles.body}>
                  <span className={styles.title}>{item.title}</span>
                  <span className={styles.meta}>{item.meta}</span>
                </span>
              </div>
            ))}
          </div>
        </section>

        <section className={`${styles.column} ${styles.wire} js-col`}>
          <div className={styles.columnHeader}>
            <span className={styles.rail} style={{ background: "var(--orange-9)" }} />
            <div className={styles.headText}>
              <span className={styles.headName}>{t("wire.desk.name")}</span>
              <span className={styles.headMeta}>{t("wire.desk.meta")}</span>
            </div>
          </div>
          <div className={styles.list}>
            {deskItems.map((item, index) => (
              <div key={item.title} className={styles.item}>
                <span className={styles.rank}>{String(index + 1).padStart(2, "0")}</span>
                <span className={styles.body}>
                  <span className={styles.title}>{item.title}</span>
                  <span className={styles.meta}>{item.meta}</span>
                </span>
              </div>
            ))}
          </div>
        </section>

        <section className={`${styles.column} ${styles.desk} js-col`}>
          <div className={styles.columnHeader}>
            <span className={styles.rail} style={{ background: "var(--accent-9)" }} />
            <div className={styles.headText}>
              <span className={styles.headName}>{t("form.column")}</span>
              <span className={styles.headMeta}>{t("form.columnMeta")}</span>
            </div>
          </div>
          <FormProvider {...form}>
            <form noValidate onSubmit={form.handleSubmit(onSubmit)} className={styles.list}>
              <Controller
                name="email"
                control={form.control}
                render={({ field, fieldState }) => (
                  <div className={`${styles.field} js-field`}>
                    <span className={styles.rank}>01</span>
                    <span className={styles.fieldBody}>
                      <label className={styles.fieldLabel} htmlFor="login-email">
                        {t("form.emailLabel")}
                      </label>
                      <TextField.Root id="login-email" size="3" type="email" autoComplete="email" placeholder={t("form.emailPlaceholder")} radius="none" {...field} />
                      {fieldState.error ? (
                        <Text size="1" color="red">
                          {fieldState.error.message}
                        </Text>
                      ) : null}
                    </span>
                  </div>
                )}
              />

              <Controller
                name="password"
                control={form.control}
                render={({ field, fieldState }) => (
                  <div className={`${styles.field} js-field`}>
                    <span className={styles.rank}>02</span>
                    <span className={styles.fieldBody}>
                      <label className={styles.fieldLabel} htmlFor="login-password">
                        {t("form.passwordLabel")}
                      </label>
                      <TextField.Root
                        id="login-password"
                        size="3"
                        type={showPassword ? "text" : "password"}
                        autoComplete="current-password"
                        placeholder={t("form.passwordPlaceholder")}
                        radius="none"
                        {...field}
                      >
                        <TextField.Slot side="right">
                          <Button
                            type="button"
                            size="1"
                            variant="ghost"
                            color="gray"
                            onClick={() => setShowPassword((v) => !v)}
                            aria-label={showPassword ? t("form.hidePassword") : t("form.showPassword")}
                          >
                            <AnimatedIcon icon={showPassword ? EyeOffIcon : EyeIcon} size={14} />
                          </Button>
                        </TextField.Slot>
                      </TextField.Root>
                      {fieldState.error ? (
                        <Text size="1" color="red">
                          {fieldState.error.message}
                        </Text>
                      ) : null}
                    </span>
                  </div>
                )}
              />

              <div className={`${styles.field} js-field`}>
                <span className={styles.rank}>03</span>
                <span className={styles.fieldBody}>
                  <Controller
                    name="rememberMe"
                    control={form.control}
                    render={({ field }) => (
                      <Flex asChild align="center" gap="2">
                        <Text as="label" size="2">
                          <Checkbox checked={!!field.value} onCheckedChange={(v) => field.onChange(v === true)} />
                          {t("form.rememberMe")}
                        </Text>
                      </Flex>
                    )}
                  />
                </span>
              </div>

              <div className={`${styles.submit} js-field`}>
                <Button type="submit" size="3" loading={login.isPending} radius="none" style={{ width: "100%" }}>
                  {t("form.submit")}
                  <AnimatedIcon icon={ArrowNarrowRightIcon} size={16} />
                </Button>
              </div>
            </form>
          </FormProvider>
        </section>

        <aside className={`${styles.column} ${styles.promo} js-col`}>
          <div className={styles.columnHeader}>
            <span className={styles.rail} style={{ background: "var(--violet-9)" }} />
            <div className={styles.headText}>
              <span className={styles.headName}>{t("promo.column")}</span>
              <span className={styles.headMeta}>{t("promo.columnMeta")}</span>
            </div>
          </div>
          <div className={styles.promoBody}>
            <p className={styles.promoBlurb}>{t("promo.blurb")}</p>
            <Link
              href={ROUTES.SIGNUP}
              size="2"
              weight="bold"
              onClick={(e) => {
                e.preventDefault();
                routerEventEmitter.navigate({ to: ROUTES.SIGNUP });
              }}
            >
              {t("promo.createAccount")}
            </Link>
          </div>
        </aside>
      </div>
    </main>
  );
}
