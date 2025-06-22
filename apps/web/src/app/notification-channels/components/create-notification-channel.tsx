import CreateEditNotificationChannel, {
  type NotificationForm,
} from "../components/create-edit-notification-channel";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getNotificationChannelsInfiniteQueryKey,
  postNotificationChannelsMutation,
} from "@/api/@tanstack/react-query.gen";
import { toast } from "sonner";
import { commonMutationErrorHandler } from "@/lib/utils";
import type {
  NotificationChannelCreateUpdateDto,
  NotificationChannelModel,
} from "@/api";

const CreateNotificationChannel = ({
  onSuccess,
}: {
  onSuccess: (notifier: NotificationChannelModel) => void;
}) => {
  const queryClient = useQueryClient();

  const createNotifierMutation = useMutation({
    ...postNotificationChannelsMutation(),
    onSuccess: (response) => {
      toast.success("Notifier created successfully");

      queryClient.invalidateQueries({
        queryKey: getNotificationChannelsInfiniteQueryKey(),
      });
      onSuccess(response.data);
    },
    onError: commonMutationErrorHandler("Failed to create notifier"),
  });

  const handleSubmit = (data: NotificationForm) => {
    const payload: NotificationChannelCreateUpdateDto = {
      name: data.name,
      type: data.type,
      config: JSON.stringify(data),
      active: true,
      is_default: false,
    };

    createNotifierMutation.mutate({
      body: payload,
    });
  };

  return <CreateEditNotificationChannel onSubmit={handleSubmit} />;
};

export default CreateNotificationChannel;
