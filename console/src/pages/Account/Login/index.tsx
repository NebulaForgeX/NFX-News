import { AnimatedIcon, ArrowNarrowRightIcon } from "nfx-ui/icons";
import type { Login } from "nfx-ui/types";

import { useMemo, useRef, useState, type ReactNode } from "react";
import { useGSAP } from "@gsap/react";
import { Box, Button, Flex, Link, Tabs, Text } from "@radix-ui/themes";
import gsap from "gsap";
import { APP_NAME } from "nfx-ui/config";
import { useLoginWithEmail, useLoginWithPhone } from "nfx-ui/hooks";
import { LoginFormData, LoginWithPhoneFormData, useInitLoginForm, useInitLoginWithPhoneForm } from "nfx-ui/schemas";
import { FormProvider, SubmitHandler } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { routerEventEmitter } from "@/events/router";
import { LoginEmailController, LoginPasswordController, LoginPhoneController, LoginRememberController } from "@/features/account";
import { ROUTES } from "@/navigations";
import AuthToolbar from "@/pages/Account/shared/AuthToolbar";
import { safeArray, safeOr } from "@/utils";

import SelectDesk from "./SelectDesk";
import styles from "./s.module.css";

gsap.registerPlugin(useGSAP);

type WireItem = { title: string; meta: string };

export default function LoginPage() {
  const { t, i18n } = useTranslation("pages.Account.Login");
  const emailForm = useInitLoginForm();
  const phoneForm = useInitLoginWithPhoneForm();
  const loginEmail = useLoginWithEmail();
  const loginPhone = useLoginWithPhone();
  const [profiles, setProfiles] = useState<Login.ProfileItem[]>([]);
  const [channel, setChannel] = useState<"email" | "phone">("email");
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

  const finishLogin = (result: Login.Response.LoginWithEmail) => {
    const list = safeArray(result?.profiles);
    if (list.length > 0) {
      setProfiles(list);
      return;
    }
    routerEventEmitter.navigate({ to: ROUTES.USER_OVERVIEW, replace: true });
  };

  const onEmail: SubmitHandler<LoginFormData> = async (data) => {
    const result = await loginEmail.mutateAsync({
      email: data.email,
      password: data.password,
      rememberMe: safeOr(data.rememberMe, false),
    });
    finishLogin(result);
  };

  const onPhone: SubmitHandler<LoginWithPhoneFormData> = async (data) => {
    const result = await loginPhone.mutateAsync({
      phone: data.phone,
      password: data.password,
      rememberMe: safeOr(data.rememberMe, false),
    });
    finishLogin(result);
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
    <Flex ref={pageRef} direction="column" className={styles.page} asChild>
      <main>
        <Box className={styles.masthead}>
          <Box className={styles.mastheadPx}>
            <Box className={styles.mastheadPy}>
              <Flex asChild align="baseline" justify="between" gap="4">
                <header>
                  <Flex className={styles.brand} align="baseline" gap="4" wrap="wrap">
                    <Text as="span" className={styles.wordmark}>
                      {APP_NAME}
                    </Text>
                    <Text as="span" className={styles.dateline}>
                      <time dateTime={new Date().toISOString()}>{dateLabel}</time>
                    </Text>
                    <Text as="span" className={styles.edition}>
                      {t("masthead.edition")}
                    </Text>
                  </Flex>
                  <AuthToolbar />
                </header>
              </Flex>
            </Box>
          </Box>
        </Box>

        <Flex className={styles.board}>
          <Box asChild className={`${styles.column} ${styles.wire} js-col`}>
            <section>
              <ColumnHead rail={styles.railWorld} name={t("wire.world.name")} meta={t("wire.world.meta")} />
              <Box className={styles.list}>
                {worldItems.map((item, index) => (
                  <WireItem key={item.title} index={index} title={item.title} meta={item.meta} />
                ))}
              </Box>
            </section>
          </Box>

          <Box asChild className={`${styles.column} ${styles.wire} js-col`}>
            <section>
              <ColumnHead rail={styles.railDesk} name={t("wire.desk.name")} meta={t("wire.desk.meta")} />
              <Box className={styles.list}>
                {deskItems.map((item, index) => (
                  <WireItem key={item.title} index={index} title={item.title} meta={item.meta} />
                ))}
              </Box>
            </section>
          </Box>

          <Box asChild className={`${styles.column} ${styles.desk} js-col`}>
            <section>
              <ColumnHead rail={styles.railForm} name={t("form.column")} meta={t("form.columnMeta")} />
              <Tabs.Root value={channel} onValueChange={(value) => setChannel(value as "email" | "phone")}>
                <Box className={styles.channelsWrap}>
                  <Box className={styles.channelsPx}>
                    <Box className={styles.channelsPy}>
                      <Tabs.List className={styles.channels}>
                        <Tabs.Trigger value="email">{t("form.channelEmail")}</Tabs.Trigger>
                        <Tabs.Trigger value="phone">{t("form.channelPhone")}</Tabs.Trigger>
                      </Tabs.List>
                    </Box>
                  </Box>
                </Box>
                <Tabs.Content value="email">
                  <FormProvider {...emailForm}>
                    <form noValidate onSubmit={emailForm.handleSubmit(onEmail)} className={styles.list}>
                      <FieldRow rank="01">
                        <LoginEmailController />
                      </FieldRow>
                      <FieldRow rank="02">
                        <LoginPasswordController />
                      </FieldRow>
                      <FieldRow rank="03">
                        <LoginRememberController />
                      </FieldRow>
                      <Box className={`${styles.submitPx} js-field`}>
                        <Box className={styles.submitPy}>
                          <Button type="submit" size="3" loading={loginEmail.isPending} radius="none" className={styles.fullWidth}>
                            {t("form.submit")}
                            <AnimatedIcon icon={ArrowNarrowRightIcon} size={16} />
                          </Button>
                        </Box>
                      </Box>
                    </form>
                  </FormProvider>
                </Tabs.Content>
                <Tabs.Content value="phone">
                  <FormProvider {...phoneForm}>
                    <form noValidate onSubmit={phoneForm.handleSubmit(onPhone)} className={styles.list}>
                      <FieldRow rank="01">
                        <LoginPhoneController />
                      </FieldRow>
                      <FieldRow rank="02">
                        <LoginPasswordController />
                      </FieldRow>
                      <FieldRow rank="03">
                        <LoginRememberController />
                      </FieldRow>
                      <Box className={`${styles.submitPx} js-field`}>
                        <Box className={styles.submitPy}>
                          <Button type="submit" size="3" loading={loginPhone.isPending} radius="none" className={styles.fullWidth}>
                            {t("form.submit")}
                            <AnimatedIcon icon={ArrowNarrowRightIcon} size={16} />
                          </Button>
                        </Box>
                      </Box>
                    </form>
                  </FormProvider>
                </Tabs.Content>
              </Tabs.Root>
            </section>
          </Box>

          <Box asChild className={`${styles.column} ${styles.promo} js-col`}>
            <aside>
              <ColumnHead rail={styles.railPromo} name={t("promo.column")} meta={t("promo.columnMeta")} />
              <Box className={styles.promoPx}>
                <Box className={styles.promoPy}>
                  <Flex direction="column" gap="3">
                    <Text as="p" className={styles.promoBlurb}>
                      {t("promo.blurb")}
                    </Text>
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
                  </Flex>
                </Box>
              </Box>
            </aside>
          </Box>
        </Flex>
      </main>
    </Flex>
  );
}

function ColumnHead({ rail, name, meta }: { rail: string; name: string; meta: string }) {
  return (
    <Box className={styles.columnHeader}>
      <Box className={styles.columnHeaderPx}>
        <Box className={styles.columnHeaderPy}>
          <Flex align="start" gap="2">
            <Box className={`${styles.rail} ${rail}`} />
            <Box className={styles.headText}>
              <Text as="span" className={styles.headName}>
                {name}
              </Text>
              <Text as="span" className={styles.headMeta}>
                {meta}
              </Text>
            </Box>
          </Flex>
        </Box>
      </Box>
    </Box>
  );
}

function WireItem({ index, title, meta }: { index: number; title: string; meta: string }) {
  return (
    <Box className={styles.item}>
      <Box className={styles.itemPx}>
        <Box className={styles.itemPy}>
          <Flex gap="2">
            <Text as="span" className={styles.rank}>
              {String(index + 1).padStart(2, "0")}
            </Text>
            <Box className={styles.body}>
              <Text as="span" className={styles.title}>
                {title}
              </Text>
              <Text as="span" className={styles.meta}>
                {meta}
              </Text>
            </Box>
          </Flex>
        </Box>
      </Box>
    </Box>
  );
}

function FieldRow({ rank, children }: { rank: string; children: ReactNode }) {
  return (
    <Box className={`${styles.field} js-field`}>
      <Box className={styles.fieldPx}>
        <Box className={styles.fieldPy}>
          <Flex gap="2">
            <Text as="span" className={styles.rank}>
              {rank}
            </Text>
            <Box className={styles.fieldBody}>{children}</Box>
          </Flex>
        </Box>
      </Box>
    </Box>
  );
}
