import { ArrowNarrowLeftIcon, FileDescriptionIcon } from "nfx-ui/icons";
import { memo } from "react";
import { Badge, Button, Flex, Link, Text } from "@radix-ui/themes";
import { useTranslation } from "react-i18next";
import { useNavigate, useParams } from "react-router";
import { DataTable, EmptyState, PageHeader, SectionBlock } from "@/components";
import { PageFrame } from "@/layouts";
import { ROUTES } from "@/navigations";
import { useDispatchReport, useOpenSnapshotHTML, useSnapshot } from "@/hooks/news";
import type { SnapshotItem } from "@/types/domain";
import { formatDateTime } from "@/utils";

const ReportDetailPage = memo(() => {
  const { t } = useTranslation("pages.Reports");
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { data, isLoading } = useSnapshot(id);
  const dispatch = useDispatchReport();
  const openHtml = useOpenSnapshotHTML();
  const items = data?.payload?.items ?? [];

  if (!isLoading && !data) {
    return (
      <PageFrame>
        <EmptyState icon={FileDescriptionIcon} title={t("emptySnapshots")} />
      </PageFrame>
    );
  }

  return (
    <PageFrame>
      <PageHeader
        icon={FileDescriptionIcon}
        title={data?.title || t("snapshots")}
        description={data?.createdAt ? formatDateTime(data.createdAt) : undefined}
        actions={
          <Flex gap="2">
            <Button variant="ghost" onClick={() => navigate(ROUTES.REPORTS)}>
              <ArrowNarrowLeftIcon size={14} />
              {t("back")}
            </Button>
            {id ? (
              <>
                <Button variant="outline" onClick={() => void openHtml.mutateAsync(id)}>
                  {t("openHtml")}
                </Button>
                <Button variant="outline" onClick={() => dispatch.mutate(id)} disabled={dispatch.isPending}>
                  {t("dispatch")}
                </Button>
              </>
            ) : null}
          </Flex>
        }
      />
      <SectionBlock
        title={t("matched")}
        description={data ? `${t("mode")}: ${data.mode} · ${t("itemCount")}: ${data.itemCount}` : undefined}
      >
        <DataTable
          loading={isLoading}
          empty={t("emptyItems")}
          rows={items}
          rowKey={(row: SnapshotItem) => row.id || `${row.sourceId}-${row.url}-${row.title}`}
          columns={[
            {
              key: "isNew",
              header: t("isNew"),
              width: "72px",
              render: (row) =>
                row.isNew ? (
                  <Badge color="amber" variant="outline">
                    {t("new")}
                  </Badge>
                ) : (
                  "—"
                ),
            },
            { key: "group", header: t("group") },
            {
              key: "title",
              header: t("titleCol"),
              render: (row) =>
                row.url ? (
                  <Link href={row.url} target="_blank" rel="noreferrer">
                    {row.title}
                  </Link>
                ) : (
                  <Text size="2">{row.title}</Text>
                ),
            },
            { key: "sourceId", header: t("source") },
          ]}
        />
      </SectionBlock>
    </PageFrame>
  );
});

ReportDetailPage.displayName = "ReportDetailPage";
export default ReportDetailPage;
