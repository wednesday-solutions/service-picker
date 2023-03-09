import { StaticSite } from "sst/constructs";

export function FrontendStack({ stack }) {
  // Deploy our React app
  const site = new StaticSite(stack, "ReactSite", {
    path: "frontend",
    buildCommand: "yarn run build",
    buildOutput: "build",
  });

  // Show the URLs in the output
  stack.addOutputs({
    SiteUrl: site.url || "http://localhost:3000/",
  });
}
