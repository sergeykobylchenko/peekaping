import HttpForm from "./http";
import PushForm from "./push";
import { useMonitorFormContext } from "../context/monitor-form-context";

const typeSpecificComponentsRegistry = {
  http: HttpForm,
  push: PushForm,
};

const CreateEditForm = () => {
  const { form } = useMonitorFormContext();
  const type = form.watch("type");

  const TypeComponent = typeSpecificComponentsRegistry[
    type as keyof typeof typeSpecificComponentsRegistry
  ] as React.ComponentType<unknown>;

  if (!TypeComponent) {
    console.log("TypeComponent not found", type);
    return null;
  }

  if (Object.keys(form.formState.errors).length > 0) {
    console.log("CreateEditFields errors", form.formState.errors);
  }

  return <TypeComponent />;
};

export default CreateEditForm;
