import Layout from "@/layout";
import { useNavigate } from "react-router-dom";
import CreateNotificationChannel from "../components/create-notification-channel";

const NewNotificationChannel = () => {
  const navigate = useNavigate();

  return (
    <Layout pageName="New Notification channel">
      <CreateNotificationChannel onSuccess={() => navigate("/notification-channels")} />
    </Layout>
  );
};

export default NewNotificationChannel;
