export const getConfig = () => {
  const config = (window as unknown as { __CONFIG__: { API_URL: string } })
    .__CONFIG__;
  const isProd = import.meta.env.PROD;

  return {
    API_URL: isProd
      ? config.API_URL // don't fallback to default in prod
      : "http://localhost:8034",
  };
};
