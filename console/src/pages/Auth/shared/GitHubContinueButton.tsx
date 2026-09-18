import { Button } from "@radix-ui/themes";
import { GithubIcon } from "nfx-ui/icons";
import { useGetGitHubAuthorizeUrl } from "nfx-ui/hooks";

export default function GitHubContinueButton({ label }: { label: string }) {
  const getUrl = useGetGitHubAuthorizeUrl();

  return (
    <Button
      type="button"
      size="3"
      variant="outline"
      color="gray"
      loading={getUrl.isPending}
      style={{ width: "100%" }}
      onClick={() => {
        void getUrl.mutateAsync().then((res) => {
          if (res?.authorizeUrl) window.location.assign(res.authorizeUrl);
        });
      }}
    >
      <GithubIcon size={16} />
      {label}
    </Button>
  );
}
