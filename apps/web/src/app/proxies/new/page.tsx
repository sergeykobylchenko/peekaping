import Layout from "@/layout";
import { useNavigate } from "react-router-dom";
import CreateProxy from "../components/create-proxy";

const NewProxy = () => {
  const navigate = useNavigate();

  return (
    <Layout pageName="New Proxy">
      <CreateProxy onSuccess={() => navigate("/proxies")} />
    </Layout>
  );
};

export default NewProxy;
