import { AnimatedIcon, ArrowNarrowRightIcon } from "nfx-ui/icons";
import type { Login } from "nfx-ui/types";

import { useMemo, useRef, useState, type ReactNode } from "react";
import { useGSAP } from "@gsap/react";
import { Section, Container, Box, Button, Flex, Link, Tabs, Text } from "@radix-ui/themes";
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
  const emailForm = useInitLoginForm(t);
  const phoneForm = useInitLoginWithPhoneForm(t);
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
    routerEventEmitter.navigate({ to: ROUTES.READER, replace: true });
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
    <Box className={`${styles.pageFill} ${styles.pageInk}`}>
    <Flex ref={pageRef} direction="column" className={styles.pageClip} asChild>
      <main>
        <Flex direction="column" flexShrink="0" className={styles.masthead}>
          <Container className={styles.mastheadPx}>
            <Section className={styles.mastheadPy}>
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
            </Section>
          </Container>
        </Flex>

        <Flex align="stretch" className={styles.board}>
          <Flex asChild direction="column" minHeight="0" align="stretch" className={`${styles.column} ${styles.columnEdge} ${styles.wire} js-col`}>
            <section>
              <ColumnHead rail={styles.railWorld} name={t("wire.world.name")} meta={t("wire.world.meta")} />
              <Flex direction="column" minHeight="0" className={styles.list}>
                {worldItems.map((item, index) => (
                  <WireItem key={item.title} index={index} title={item.title} meta={item.meta} />
                ))}
              </Flex>
            </section>
          </Flex>

          <Flex asChild direction="column" minHeight="0" align="stretch" className={`${styles.column} ${styles.columnEdge} ${styles.wire} js-col`}>
            <section>
              <ColumnHead rail={styles.railDesk} name={t("wire.desk.name")} meta={t("wire.desk.meta")} />
              <Flex direction="column" minHeight="0" className={styles.list}>
                {deskItems.map((item, index) => (
                  <WireItem key={item.title} index={index} title={item.title} meta={item.meta} />
                ))}
              </Flex>
            </section>
          </Flex>

          <Flex asChild direction="column" minHeight="0" align="stretch" className={`${styles.column} ${styles.columnEdge} ${styles.desk} js-col`}>
            <section>
              <ColumnHead rail={styles.railForm} name={t("form.column")} meta={t("form.columnMeta")} />
              <Tabs.Root value={channel} onValueChange={(value) => setChannel(value as "email" | "phone")}>
                <Box className={styles.channelsWrap}>
                  <Container className={styles.channelsPx}>
                    <Section className={styles.channelsPy}>
                      <Tabs.List className={`${styles.channelsSize} ${styles.channelsShadow} ${styles.channelsFill}`}>
                        <Tabs.Trigger value="email">{t("form.channelEmail")}</Tabs.Trigger>
                        <Tabs.Trigger value="phone">{t("form.channelPhone")}</Tabs.Trigger>
                      </Tabs.List>
                    </Section>
                  </Container>
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
                      <Container className={`${styles.submitPx} js-field`}>
                        <Section className={styles.submitPy}>
                          <Button type="submit" size="3" loading={loginEmail.isPending} radius="none" className={styles.fullWidth}>
                            {t("form.submit")}
                            <AnimatedIcon icon={ArrowNarrowRightIcon} size={16} />
                          </Button>
                        </Section>
                      </Container>
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
                      <Container className={`${styles.submitPx} js-field`}>
                        <Section className={styles.submitPy}>
                          <Button type="submit" size="3" loading={loginPhone.isPending} radius="none" className={styles.fullWidth}>
                            {t("form.submit")}
                            <AnimatedIcon icon={ArrowNarrowRightIcon} size={16} />
                          </Button>
                        </Section>
                      </Container>
                    </form>
                  </FormProvider>
                </Tabs.Content>
              </Tabs.Root>
            </section>
          </Flex>

          <Flex asChild direction="column" minHeight="0" align="stretch" className={`${styles.column} ${styles.columnEdge} ${styles.promo} js-col`}>
            <aside>
              <ColumnHead rail={styles.railPromo} name={t("promo.column")} meta={t("promo.columnMeta")} />
              <Container className={styles.promoPx}>
                <Section className={styles.promoPy}>
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
                </Section>
              </Container>
            </aside>
          </Flex>
        </Flex>
      </main>
    </Flex>
    </Box>
  );
}

function ColumnHead({ rail, name, meta }: { rail: string; name: string; meta: string }) {
  return (
    <Flex direction="column" flexShrink="0" className={styles.columnHeader}>
      <Container className={styles.columnHeaderPx}>
        <Section className={styles.columnHeaderPy}>
          <Flex align="start" gap="2">
            <Section className={`${styles.railItem} ${styles.railOffset}`}>
              <Box className={`${styles.railSize} ${styles.railRadius} ${rail}`} />
            </Section>
            <Flex direction="column" gap="1" minWidth="0" flexGrow="1" className={styles.headText}>
              <Text as="span" className={styles.headName}>
                {name}
              </Text>
              <Section className={styles.headMetaSpace}>
                <Text as="span" className={styles.headMeta}>
                  {meta}
                </Text>
              </Section>
            </Flex>
          </Flex>
        </Section>
      </Container>
    </Flex>
  );
}

function WireItem({ index, title, meta }: { index: number; title: string; meta: string }) {
  return (
    <Box className={styles.item}>
      <Container className={styles.itemPx}>
        <Section className={styles.itemPy}>
          <Flex gap="2">
            <Flex className={styles.rank}>
              <Text as="span">{String(index + 1).padStart(2, "0")}</Text>
            </Flex>
            <Flex direction="column" gap="1" minWidth="0" className={styles.body}>
              <Text as="span" className={styles.title}>
                {title}
              </Text>
              <Text as="span" className={styles.meta}>
                {meta}
              </Text>
            </Flex>
          </Flex>
        </Section>
      </Container>
    </Box>
  );
}

function FieldRow({ rank, children }: { rank: string; children: ReactNode }) {
  return (
    <Box className={`${styles.field} js-field`}>
      <Container className={styles.fieldPx}>
        <Section className={styles.fieldPy}>
          <Flex gap="2">
            <Flex className={styles.rank}>
              <Text as="span">{rank}</Text>
            </Flex>
            <Flex minWidth="0" flexGrow="1" className={styles.fieldBody}>{children}</Flex>
          </Flex>
        </Section>
      </Container>
    </Box>
  );
}
