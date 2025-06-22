import {
  getProxiesInfiniteQueryKey,
  postProxiesMutation,
} from "@/api/@tanstack/react-query.gen";
import { commonMutationErrorHandler } from "@/lib/utils";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import type { ProxyForm } from "./create-edit-proxy";
import type { ProxyCreateUpdateDto, ProxyModel } from "@/api";
import CreateEditProxy from "./create-edit-proxy";

const CreateProxy = ({
  onSuccess,
}: {
  onSuccess: (proxy: ProxyModel) => void;
}) => {
  const queryClient = useQueryClient();

  const createProxyMutation = useMutation({
    ...postProxiesMutation(),
    onSuccess: (response) => {
      toast.success("Proxy created successfully");
      queryClient.invalidateQueries({ queryKey: getProxiesInfiniteQueryKey() });
      onSuccess(response.data);
    },
    onError: commonMutationErrorHandler("Failed to create proxy"),
  });

  const handleSubmit = (data: ProxyForm) => {
    const proxyData: ProxyCreateUpdateDto = {
      protocol: data.protocol,
      host: data.host,
      port: data.port,
      auth: data.auth,
      username: data.auth ? data.username : undefined,
      password: data.auth ? data.password : undefined,
    };

    createProxyMutation.mutate({
      body: proxyData,
    });
  };

  return <CreateEditProxy onSubmit={handleSubmit} />;
};

export default CreateProxy;
