import { AnimatedIcon, ArrowNarrowRightIcon } from "nfx-ui/icons";
import { useRef } from "react";
import { useGSAP } from "@gsap/react";
import { Box, Button, Flex, Heading, Link, Text } from "@radix-ui/themes";
import gsap from "gsap";
import { APP_NAME } from "nfx-ui/config";
import { AuthSignupPlatformEnum, LanguageEnum } from "nfx-ui/enums";
import { useSendVerificationCode, useSignupWithEmail } from "nfx-ui/hooks";
import { SignupFormData, useInitSignupForm } from "nfx-ui/schemas";
import { usePreferenceStore } from "nfx-ui/stores";
import { FormProvider, SubmitHandler } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { routerEventEmitter } from "@/events/router";
import {
  SignupConfirmPasswordController,
  SignupEmailController,
  SignupPasswordController,
  SignupRememberController,
  SignupVerificationCodeController,
} from "@/features/account";
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
      <Box className={styles.ticker}>
        <Box className={styles.tickerPx}>
          <Box className={styles.tickerPy}>
            <Flex align="center" justify="between" gap="4">
              <span className={styles.tickerText}>{t("ticker")}</span>
              <AuthToolbar />
            </Flex>
          </Box>
        </Box>
      </Box>

      <div className={styles.body}>
        <Box asChild className={`${styles.ticket} js-ticket`}>
          <section>
            <Box className={styles.ticketPx}>
              <Box className={styles.ticketPy}>
                <Text as="span" className={styles.kicker}>
                  {t("ticket.kicker")}
                </Text>
                <Heading as="h1" size="7" className={styles.headline}>
                  {t("ticket.title", { name: APP_NAME })}
                </Heading>
                <Text as="p" className={styles.lede}>
                  {t("ticket.subtitle")}
                </Text>

                <FormProvider {...form}>
                  <Flex asChild direction="column" gap="4">
                    <form noValidate onSubmit={form.handleSubmit(onSubmit)} className={styles.form}>
                      <SignupEmailController helperText={t("emailHint")} />
                      <Flex direction="column" gap="2">
                        <SignupVerificationCodeController />
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
                      </Flex>
                      <SignupPasswordController />
                      <SignupConfirmPasswordController />
                      <SignupRememberController />

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
              </Box>
            </Box>
          </section>
        </Box>

        <Box asChild className={`${styles.briefing} js-brief`}>
          <aside>
            <Box className={styles.briefingPx}>
              <Box className={styles.briefingPy}>
                <span className={styles.briefKicker}>{t("briefing.kicker")}</span>
                <h2 className={styles.briefTitle}>{t("briefing.title")}</h2>
                {steps.map((step) => (
                  <Box key={step.rank} className={`${styles.step} js-step`}>
                    <Box className={styles.stepPy}>
                      <div className={styles.stepGrid}>
                        <span className={styles.stepRank}>{step.rank}</span>
                        <span>
                          <span className={styles.stepTitle}>{step.title}</span>
                          <span className={styles.stepBody}>{step.body}</span>
                        </span>
                      </div>
                    </Box>
                  </Box>
                ))}
              </Box>
            </Box>
          </aside>
        </Box>
      </div>
    </main>
  );
}
