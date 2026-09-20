import { AnimatedIcon, ArrowNarrowRightIcon, EyeIcon, EyeOffIcon } from "nfx-ui/icons";
import { useRef, useState } from "react";
import { useGSAP } from "@gsap/react";
import { Button, Checkbox, Flex, Link, Text, TextField } from "@radix-ui/themes";
import gsap from "gsap";
import { APP_NAME } from "nfx-ui/config";
import { AuthSignupPlatformEnum, LanguageEnum } from "nfx-ui/enums";
import { useSendVerificationCode, useSignupWithEmail } from "nfx-ui/hooks";
import { SignupFormData, useInitSignupForm } from "nfx-ui/schemas";
import { usePreferenceStore } from "nfx-ui/stores";
import { Controller, FormProvider, SubmitHandler } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { routerEventEmitter } from "@/events/router";
import { ROUTES } from "@/navigations";
import AuthToolbar from "@/pages/Account/shared/AuthToolbar";

import styles from "./s.module.css";

gsap.registerPlugin(useGSAP);

type BriefStep = { rank: string; title: string; body: string };

export default function SignupPage() {
  const { t } = useTranslation("pages.Account.Signup");
  const form = useInitSignupForm();
  const signup = useSignupWithEmail();
  const sendCode = useSendVerificationCode();
  const language = usePreferenceStore((s) => s.language);
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const pageRef = useRef<HTMLDivElement>(null);

  useGSAP(
    () => {
      if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
      gsap.set(".js-ticket", { autoAlpha: 0, x: -40 });
      gsap.set(".js-brief", { autoAlpha: 0, x: 40 });
      gsap.set(".js-step", { autoAlpha: 0, y: 12 });
      const tl = gsap.timeline({ defaults: { ease: "power3.out" } });
      tl.to(".js-ticket", { autoAlpha: 1, x: 0, duration: 0.55 })
        .to(".js-brief", { autoAlpha: 1, x: 0, duration: 0.55 }, "-=0.35")
        .to(".js-step", { autoAlpha: 1, y: 0, duration: 0.4, stagger: 0.08 }, "-=0.3");
    },
    { scope: pageRef },
  );

  const onSubmit: SubmitHandler<SignupFormData> = async (data) => {
    await signup.mutateAsync({
      email: data.email,
      password: data.password,
      verificationCode: data.verificationCode,
      lang: language ?? LanguageEnum.EN,
      rememberMe: data.rememberMe ?? false,
      signupPlatform: AuthSignupPlatformEnum.NFXNEWS,
    });
    routerEventEmitter.navigate({ to: ROUTES.USER_OVERVIEW, replace: true });
  };

  const email = form.watch("email");
  const stepsRaw = t("briefing.steps", { returnObjects: true });
  const steps = Array.isArray(stepsRaw) ? (stepsRaw as BriefStep[]) : [];

  return (
    <main ref={pageRef} className={styles.page}>
      <div className={styles.ticker}>
        <span className={styles.tickerText}>{t("ticker")}</span>
        <AuthToolbar />
      </div>

      <div className={styles.body}>
        <section className={`${styles.ticket} js-ticket`}>
          <span className={styles.kicker}>{t("ticket.kicker")}</span>
          <h1 className={styles.headline}>{t("ticket.title", { name: APP_NAME })}</h1>
          <p className={styles.lede}>{t("ticket.subtitle")}</p>

          <FormProvider {...form}>
            <Flex asChild direction="column" gap="4">
              <form noValidate onSubmit={form.handleSubmit(onSubmit)} className={styles.form}>
                <Controller
                  name="email"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Flex direction="column" gap="1">
                      <Text as="label" size="2" weight="medium" htmlFor="signup-email">
                        {t("emailLabel")}
                      </Text>
                      <TextField.Root id="signup-email" size="3" type="email" autoComplete="email" placeholder={t("emailPlaceholder")} radius="none" {...field} />
                      <span className={styles.hint}>{t("emailHint")}</span>
                      {fieldState.error ? (
                        <Text size="1" color="red">
                          {fieldState.error.message}
                        </Text>
                      ) : null}
                    </Flex>
                  )}
                />

                <Flex direction="column" gap="1">
                  <Text as="label" size="2" weight="medium" htmlFor="signup-code">
                    {t("codeLabel")}
                  </Text>
                  <div className={styles.wireRow}>
                    <Controller
                      name="verificationCode"
                      control={form.control}
                      render={({ field }) => <TextField.Root id="signup-code" size="3" placeholder={t("codePlaceholder")} radius="none" {...field} />}
                    />
                    <Button
                      type="button"
                      size="3"
                      variant="outline"
                      radius="none"
                      loading={sendCode.isPending}
                      disabled={!email}
                      onClick={() =>
                        email &&
                        sendCode.mutate({
                          email,
                          lang: language ?? LanguageEnum.EN,
                        })
                      }
                    >
                      {t("sendCode")}
                    </Button>
                  </div>
                </Flex>

                <Controller
                  name="password"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Flex direction="column" gap="1">
                      <Text as="label" size="2" weight="medium" htmlFor="signup-password">
                        {t("passwordLabel")}
                      </Text>
                      <TextField.Root
                        id="signup-password"
                        size="3"
                        type={showPassword ? "text" : "password"}
                        autoComplete="new-password"
                        placeholder={t("passwordPlaceholder")}
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
                            aria-label={showPassword ? t("hidePassword") : t("showPassword")}
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
                    </Flex>
                  )}
                />

                <Controller
                  name="confirmPassword"
                  control={form.control}
                  render={({ field, fieldState }) => (
                    <Flex direction="column" gap="1">
                      <Text as="label" size="2" weight="medium" htmlFor="signup-confirm">
                        {t("confirmLabel")}
                      </Text>
                      <TextField.Root
                        id="signup-confirm"
                        size="3"
                        type={showConfirm ? "text" : "password"}
                        autoComplete="new-password"
                        placeholder={t("confirmPlaceholder")}
                        radius="none"
                        {...field}
                      >
                        <TextField.Slot side="right">
                          <Button
                            type="button"
                            size="1"
                            variant="ghost"
                            color="gray"
                            onClick={() => setShowConfirm((v) => !v)}
                            aria-label={showConfirm ? t("hidePassword") : t("showPassword")}
                          >
                            <AnimatedIcon icon={showConfirm ? EyeOffIcon : EyeIcon} size={14} />
                          </Button>
                        </TextField.Slot>
                      </TextField.Root>
                      {fieldState.error ? (
                        <Text size="1" color="red">
                          {fieldState.error.message}
                        </Text>
                      ) : null}
                    </Flex>
                  )}
                />

                <Controller
                  name="rememberMe"
                  control={form.control}
                  render={({ field }) => (
                    <Flex asChild align="center" gap="2">
                      <Text as="label" size="2">
                        <Checkbox checked={!!field.value} onCheckedChange={(v) => field.onChange(v === true)} />
                        {t("rememberMe")}
                      </Text>
                    </Flex>
                  )}
                />

                <Button type="submit" size="3" loading={signup.isPending} radius="none" style={{ width: "100%" }}>
                  {t("submit")}
                  <AnimatedIcon icon={ArrowNarrowRightIcon} size={16} />
                </Button>

                <p className={styles.signIn}>
                  {t("hasAccount")}{" "}
                  <Link
                    href={ROUTES.LOGIN}
                    size="2"
                    weight="bold"
                    onClick={(e) => {
                      e.preventDefault();
                      routerEventEmitter.navigate({ to: ROUTES.LOGIN });
                    }}
                  >
                    {t("signIn")}
                  </Link>
                </p>
              </form>
            </Flex>
          </FormProvider>
        </section>

        <aside className={`${styles.briefing} js-brief`}>
          <span className={styles.briefKicker}>{t("briefing.kicker")}</span>
          <h2 className={styles.briefTitle}>{t("briefing.title")}</h2>
          {steps.map((step) => (
            <div key={step.rank} className={`${styles.step} js-step`}>
              <span className={styles.stepRank}>{step.rank}</span>
              <span>
                <span className={styles.stepTitle}>{step.title}</span>
                <span className={styles.stepBody}>{step.body}</span>
              </span>
            </div>
          ))}
        </aside>
      </div>
    </main>
  );
}
