import type { ReactNode } from "react";

import { Box, Section, Flex, Heading, Text } from "@radix-ui/themes";

type SectionBlockProps = {
  title: string;
  description?: string;
  actions?: ReactNode;
  children: ReactNode;
};

export function SectionBlock({ title, description, actions, children }: SectionBlockProps) {
  return (
    <Flex direction="column" gap="3">
      <Flex align="end" justify="between" gap="3" wrap="wrap">
        <Flex direction="column" gap="1" minWidth="0">
          <Flex align="center" gap="2">
            <Box style={{ width: 3, height: 16, background: "var(--accent-9)", borderRadius: "var(--radius-1)" }} />
            <Heading as="h2" size="3" style={{ fontFamily: "var(--heading-font-family)", letterSpacing: "-0.02em" }}>
              {title}
            </Heading>
          </Flex>
          {description ? (
            <Section mt="1">
              <Text as="p" size="1" color="gray">
                {description}
              </Text>
            </Section>
          ) : null}
        </Flex>
        {actions ? (
          <Flex gap="2" wrap="wrap" align="center">
            {actions}
          </Flex>
        ) : null}
      </Flex>
      <Box style={{ borderTop: "1px solid var(--gray-a5)" }}>
        <Section pt="3">
          {children}
        </Section>
      </Box>
    </Flex>
  );
}
