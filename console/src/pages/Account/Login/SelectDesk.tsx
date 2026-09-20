import { AnimatedIcon, RightChevron, ShieldCheck, UsersIcon } from "nfx-ui/icons";
import type { Login } from "nfx-ui/types";

import { useMemo, useRef } from "react";
import { useGSAP } from "@gsap/react";
import { Avatar, Badge, Button, Flex, Spinner } from "@radix-ui/themes";
import gsap from "gsap";
import { ProfileKindEnum } from "nfx-ui/enums";
import { useSelectProfile } from "nfx-ui/hooks";
import { useTranslation } from "react-i18next";

import { routerEventEmitter } from "@/events/router";
import { ROUTES } from "@/navigations";
import AuthToolbar from "@/pages/Account/shared/AuthToolbar";
import { buildImageUrl, resolveAccountDisplayName, resolveAccountInitial, safeArray, safeStringable } from "@/utils";

import styles from "./SelectDesk.module.css";

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
    const forger = profiles.filter((p) => p.kind === ProfileKindEnum.FORGER);
    const authority = profiles.filter((p) => p.kind === ProfileKindEnum.AUTHORITY);
    return [
      { kind: ProfileKindEnum.FORGER as const, items: forger, rail: "var(--accent-9)" },
      { kind: ProfileKindEnum.AUTHORITY as const, items: authority, rail: "var(--amber-9)" },
    ].filter((group) => group.items.length > 0);
  }, [profiles]);

  return (
    <div ref={pageRef} className={styles.page}>
      <header className={styles.masthead}>
        <div className={styles.lead}>
          <span className={styles.kicker}>{t("selectProfile.eyebrow")}</span>
          <h1 className={styles.title}>{t("selectProfile.title")}</h1>
          <p className={styles.subtitle}>{t("selectProfile.subtitle")}</p>
        </div>
        <div className={styles.actions}>
          <Button type="button" variant="soft" color="gray" size="2" onClick={onBack} disabled={selectProfile.isPending}>
            {t("selectProfile.back")}
          </Button>
          <AuthToolbar />
        </div>
      </header>

      <div className={styles.board}>
        {groups.map((group) => {
          const isAuthority = group.kind === ProfileKindEnum.AUTHORITY;
          const kindLabel = t(`selectProfile.kind.${group.kind}`);
          return (
            <section key={group.kind} className={`${styles.column} js-desk-col`}>
              <div className={styles.columnHeader}>
                <span className={styles.rail} style={{ background: group.rail }} />
                <div>
                  <Flex align="center" gap="2">
                    <AnimatedIcon icon={isAuthority ? ShieldCheck : UsersIcon} size={14} />
                    <span className={styles.headName}>{kindLabel}</span>
                  </Flex>
                  <span className={styles.headMeta}>{t(`selectProfile.kindHint.${group.kind}`)}</span>
                </div>
              </div>
              <div className={styles.list}>
                {group.items.map((profile, index) => {
                  const name = resolveAccountDisplayName(profile.displayName, profile.profileId);
                  const initial = resolveAccountInitial(profile.displayName, profile.profileId);
                  const roles = safeArray(profile.roles);
                  const place = [safeStringable(profile.city), safeStringable(profile.country)].filter(Boolean).join(", ");
                  return (
                    <button
                      key={`${profile.kind}:${profile.profileId}`}
                      type="button"
                      className={styles.item}
                      disabled={selectProfile.isPending}
                      onClick={async () => {
                        await selectProfile.mutateAsync({
                          profileId: profile.profileId,
                          kind: profile.kind ?? ProfileKindEnum.FORGER,
                        });
                        routerEventEmitter.navigate({
                          to: ROUTES.USER_OVERVIEW,
                          replace: true,
                        });
                      }}
                    >
                      <span className={styles.rank}>{String(index + 1).padStart(2, "0")}</span>
                      <Avatar size="2" radius="none" fallback={initial} src={profile.avatarImageId ? buildImageUrl(profile.avatarImageId) : undefined} />
                      <span className={styles.body}>
                        <span className={styles.name}>{name}</span>
                        {place ? <span className={styles.place}>{place}</span> : null}
                        {roles.length > 0 ? (
                          <span className={styles.roles}>
                            {roles.map((role) => (
                              <Badge key={role} variant="soft" size="1">
                                {t(`selectProfile.roles.${role}`, { defaultValue: role })}
                              </Badge>
                            ))}
                          </span>
                        ) : null}
                      </span>
                      <span className={styles.chevron}>
                        {selectProfile.isPending ? <Spinner size="1" /> : <AnimatedIcon icon={RightChevron} size={16} />}
                      </span>
                    </button>
                  );
                })}
              </div>
            </section>
          );
        })}
      </div>
    </div>
  );
}
