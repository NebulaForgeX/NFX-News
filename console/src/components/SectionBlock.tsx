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
        <Box minWidth="0">
          <Heading as="h2" size="3">
            {title}
          </Heading>
          {description ? (
            <Section mt="1">
              <Text as="p" size="1" color="gray">
                {description}
              </Text>
            </Section>
          ) : null}
        </Box>
        {actions ? (
          <Flex gap="2" wrap="wrap" align="center">
            {actions}
          </Flex>
        ) : null}
      </Flex>
      <Box style={{ borderTop: "1px solid var(--gray-a4)" }}>
        <Section pt="3">
          {children}
        </Section>
      </Box>
    </Flex>
  );
}
