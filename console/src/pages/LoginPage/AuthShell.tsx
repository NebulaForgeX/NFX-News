import type { ReactNode } from "react";

import { useRef } from "react";
import { useGSAP } from "@gsap/react";
import { Box } from "@radix-ui/themes";
import gsap from "gsap";
import { PreferencesPopover } from "nfx-ui/components";

import styles from "./AuthShell.module.css";

gsap.registerPlugin(useGSAP);

export type AuthShellProps = {
  brandTitle: string;
  brandEyebrow?: string;
  heroFooter: string;
  children: ReactNode;
};

const BRIEFS = [
  { title: "Sources keep their own masthead.", body: "Catalog stays global; your reading prefs stay on the account." },
  { title: "Edition, not a product tour.", body: "Crawl, report, and notify sit behind this subscription slip." },
  { title: "Wire copy, not a glass logo.", body: "Identity still issues the token. The desk is the paper." },
] as const;

export default function AuthShell({ brandTitle, brandEyebrow, heroFooter, children }: AuthShellProps) {
  const pageRef = useRef<HTMLDivElement>(null);

  useGSAP(
    () => {
      if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
      gsap.set(".js-mast", { autoAlpha: 0, y: -24 });
      gsap.set(".js-brief", { autoAlpha: 0, y: 16 });
      gsap.set(".js-slip", { autoAlpha: 0, x: 28 });
      const tl = gsap.timeline({ defaults: { ease: "power2.out" } });
      tl.to(".js-mast", { autoAlpha: 1, y: 0, duration: 0.6 })
        .to(".js-brief", { autoAlpha: 1, y: 0, duration: 0.45, stagger: 0.1 }, "-=0.2")
        .to(".js-slip", { autoAlpha: 1, x: 0, duration: 0.55 }, "-=0.35");
    },
    { scope: pageRef },
  );

  return (
    <Box ref={pageRef} className={styles.page} asChild>
      <main className={styles.sheet}>
        <header className={`${styles.masthead} js-mast`}>
          <div className={styles.kicker}>
            <span>{brandEyebrow ?? "NFX News"}</span>
            <span>Vol. 01 · City edition</span>
          </div>
          <h1 className={styles.mastTitle}>The Record</h1>
          <div className={styles.rule}>Late wire</div>
        </header>
        <div className={styles.body}>
          <section className={styles.column}>
            <h2 className={styles.headline}>{brandTitle}</h2>
            <p className={styles.lede}>{heroFooter}</p>
            <div className={styles.briefs}>
              {BRIEFS.map((brief) => (
                <article key={brief.title} className={`${styles.brief} js-brief`}>
                  <h2>{brief.title}</h2>
                  <p>{brief.body}</p>
                </article>
              ))}
            </div>
          </section>
          <aside className={`${styles.slip} js-slip`}>
            <div className={styles.toolbar}>
              <PreferencesPopover />
            </div>
            <div className={styles.slipInner}>{children}</div>
          </aside>
        </div>
      </main>
    </Box>
  );
}
