import { AnimatedIcon, RightChevron, ShieldCheck, UsersIcon } from "nfx-ui/icons";
import type { Login } from "nfx-ui/types";

import { useMemo, useRef } from "react";
import { useGSAP } from "@gsap/react";
import { Section, Container, Avatar, Badge, Box, Button, Flex, Heading, Spinner, Text } from "@radix-ui/themes";
import gsap from "gsap";
import { ProfileKind, ProfileKindEnum } from "nfx-ui/enums";
import { useSelectProfile } from "nfx-ui/hooks";
import { useTranslation } from "react-i18next";

import { routerEventEmitter } from "@/events/router";
import { ROUTES } from "@/navigations";
import AuthToolbar from "@/pages/Account/shared/AuthToolbar";
import { buildImageUrl, resolveAccountDisplayName, resolveAccountInitial, safeArray, safeStringable } from "@/utils";

import styles from "./s.module.css";

gsap.registerPlugin(useGSAP);

export type SelectDeskProps = {
  profiles: Login.ProfileItem[];
  onBack: () => void;
};

export default function SelectDesk({ profiles, onBack }: SelectDeskProps) {
  const { t } = useTranslation("pages.Account.Login");
  const selectProfile = useSelectProfile();
  const pageRef = useRef<HTMLDivElement>(null);

  useGSAP(
    () => {
      if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
      gsap.set(".js-desk-col", { autoAlpha: 0, x: -28 });
      gsap.to(".js-desk-col", { autoAlpha: 1, x: 0, duration: 0.5, stagger: 0.1, ease: "power3.out" });
    },
    { scope: pageRef },
  );

  const groups = useMemo(() => {
    const community = profiles.filter((p) => ProfileKind(p.kind) === ProfileKindEnum.COMMUNITY);
    const authority = profiles.filter((p) => ProfileKind(p.kind) === ProfileKindEnum.AUTHORITY);
    return [
      { kind: ProfileKindEnum.COMMUNITY, items: community, rail: styles.railCommunity },
      { kind: ProfileKindEnum.AUTHORITY, items: authority, rail: styles.railAuthority },
    ].filter((group) => group.items.length > 0);
  }, [profiles]);

  return (
    <Flex ref={pageRef} direction="column" className={`${styles.pageClip} ${styles.pageFill} ${styles.pageInk}`}>
        <Flex direction="column" flexShrink="0" className={styles.masthead}>
          <Container className={styles.mastheadPx}>
            <Section className={styles.mastheadPy}>
              <Flex asChild align="start" justify="between" gap="4">
                <header>
                  <Flex direction="column" gap="1" className={styles.lead}>
                    <Text as="span" className={styles.kicker}>
                      {t("selectProfile.eyebrow")}
                    </Text>
                    <Heading as="h1" size="6" className={styles.title}>
                      {t("selectProfile.title")}
                    </Heading>
                    <Text as="p" className={styles.subtitle}>
                      {t("selectProfile.subtitle")}
                    </Text>
                  </Flex>
                  <Flex className={styles.actions} align="center" gap="2">
                    <Button type="button" variant="outline" color="gray" size="2" onClick={onBack} disabled={selectProfile.isPending}>
                      {t("selectProfile.back")}
                    </Button>
                    <AuthToolbar />
                  </Flex>
                </header>
              </Flex>
            </Section>
          </Container>
        </Flex>

        <Flex align="stretch" className={styles.board}>
          {groups.map((group) => {
            const isAuthority = group.kind === ProfileKindEnum.AUTHORITY;
            const kindLabel = t(`selectProfile.kind.${group.kind}`);
            return (
              <Flex asChild key={group.kind} direction="column" minHeight="0" align="stretch" className={`${styles.column} ${styles.columnSize} ${styles.columnEdge} js-desk-col`}>
                <section>
                  <Flex direction="column" flexShrink="0" className={styles.columnHeader}>
                    <Container className={styles.columnHeaderPx}>
                      <Section className={styles.columnHeaderPy}>
                        <Flex align="start" gap="2">
                          <Section className={`${styles.railItem} ${styles.railOffset}`}>
                            <Box className={`${styles.railSize} ${styles.railRadius} ${group.rail}`} />
                          </Section>
                          <Box>
                            <Flex align="center" gap="2">
                              <AnimatedIcon icon={isAuthority ? ShieldCheck : UsersIcon} size={14} />
                              <Text as="span" className={styles.headName}>
                                {kindLabel}
                              </Text>
                            </Flex>
                            <Section className={styles.headMetaSpace}>
                            <Text as="span" className={styles.headMeta}>
                              {t(`selectProfile.kindHint.${group.kind}`)}
                            </Text>
                            </Section>
                          </Box>
                        </Flex>
                    </Section>
                    </Container>
                  </Flex>
                  <Flex direction="column" minHeight="0" className={styles.list}>
                    {group.items.map((profile, index) => {
                      const name = resolveAccountDisplayName(profile.displayName, profile.profileId);
                      const initial = resolveAccountInitial(profile.displayName, profile.profileId);
                      const roles = safeArray(profile.roles);
                      const place = [safeStringable(profile.city), safeStringable(profile.country)].filter(Boolean).join(", ");
                      const kind = ProfileKind(profile.kind);
                      return (
                        <Button
                          key={`${kind}:${profile.profileId}`}
                          type="button"
                          variant="ghost"
                          className={`${styles.itemHit} ${styles.itemPadX} ${styles.itemPadY} ${styles.itemFill}`}
                          disabled={selectProfile.isPending}
                          onClick={async () => {
                            await selectProfile.mutateAsync({
                              profileId: profile.profileId,
                              kind,
                            });
                            routerEventEmitter.navigate({
                              to: ROUTES.USER_OVERVIEW,
                              replace: true,
                            });
                          }}
                        >
                          <Box className={`${styles.itemRule} ${styles.itemWide}`}>
                          <Container className={styles.itemPx}>
                            <Section className={styles.itemPy}>
                              <Flex align="start" gap="2" width="100%">
                                <Flex className={styles.rank}>
                                  <Text as="span">{String(index + 1).padStart(2, "0")}</Text>
                                </Flex>
                                <Avatar size="2" radius="none" fallback={initial} src={profile.avatarImageId ? buildImageUrl(profile.avatarImageId) : undefined} />
                                <Flex direction="column" gap="1" className={styles.body}>
                                  <Text as="span" className={styles.name}>
                                    {name}
                                  </Text>
                                  {place ? (
                                    <Text as="span" className={styles.place}>
                                      {place}
                                    </Text>
                                  ) : null}
                                  {roles.length > 0 ? (
                                    <Flex className={styles.roles} gap="1" wrap="wrap">
                                      {roles.map((role) => (
                                        <Badge key={role} variant="outline" size="1">
                                          {t(`selectProfile.roles.${role}`, { defaultValue: role })}
                                        </Badge>
                                      ))}
                                    </Flex>
                                  ) : null}
                                </Flex>
                                <Flex align="center" className={styles.chevron}>
                                  {selectProfile.isPending ? <Spinner size="1" /> : <AnimatedIcon icon={RightChevron} size={16} />}
                                </Flex>
                              </Flex>
                            </Section>
                          </Container>
                          </Box>
                        </Button>
                      );
                    })}
                  </Flex>
                </section>
              </Flex>
            );
          })}
        </Flex>
    </Flex>
  );
}
