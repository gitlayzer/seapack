// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";
import tailwindcss from "@tailwindcss/vite";
import starlightLlmsTxt from "starlight-llms-txt";
import starlightPageActions from "starlight-page-actions";

// https://astro.build/config
export default defineConfig({
  site: "https://seapack.sealos.io",

  prefetch: {
    prefetchAll: true,
    defaultStrategy: "hover",
  },

  integrations: [
    starlight({
      title: "SeaPack Docs",
      social: [
        {
          icon: "github",
          label: "GitHub",
          href: "https://github.com/gitlayzer/seapack",
        },
      ],
      editLink: {
        baseUrl: "https://github.com/gitlayzer/seapack/edit/main/docs/",
      },
      favicon: "/favicon.svg?v=2",
      customCss: [
        "./src/tailwind.css",

        "@fontsource/inter/400.css",
        "@fontsource/inter/600.css",
      ],
      plugins: [
        starlightPageActions(),
        starlightLlmsTxt({
          projectName: "SeaPack",
          description:
            "Sealos-oriented application builder that analyzes Node, Python, Go, Java, and Deno projects and turns them into container images.",
          details:
            "SeaPack provides a seamless way to build container images from your source code without complex configuration. It automatically detects your project type and generates appropriate build steps.",
          customSets: [
            {
              label: "Languages Reference",
              description:
                "Language-specific documentation for SeaPack's supported Sealos workloads",
              paths: ["languages/**"],
            },
            {
              label: "Architecture",
              description:
                "Technical details about SeaPack's internal architecture",
              paths: ["architecture/**"],
            },
            {
              label: "Guides",
              description: "Step-by-step guides for common tasks",
              paths: ["guides/**"],
            },
            {
              label: "Configuration",
              description: "Configuration options and environment variables",
              paths: ["config/**"],
            },
            {
              label: "Deploying",
              description: "Deployment guides for Sealos and GitHub Actions",
              paths: ["deploying/**"],
            },
            {
              label: "Reference",
              description: "CLI commands and BuildKit frontend reference",
              paths: ["reference/**"],
            },
          ],
          optionalLinks: [
            {
              label: "SeaPack GitHub Repository",
              url: "https://github.com/gitlayzer/seapack",
              description: "Source code and issue tracking for SeaPack",
            },
            {
              label: "Sealos",
              url: "https://sealos.com",
              description: "Cloud platform targeted by SeaPack",
            },
            {
              label: "Sealos SeaPack Guide",
              url: "https://docs.sealos.com/guides/build-configuration#seapack",
              description: "How to use SeaPack on Sealos platform",
            },
          ],
          promote: ["index*", "getting-started*", "installation*", "config/**"],
        }),
      ],
      sidebar: [
        {
          label: "Getting Started",
          link: "/getting-started",
        },
        {
          label: "FAQ",
          link: "/faq",
        },
        {
          label: "Installation",
          link: "/installation",
        },
        {
          label: "Help",
          link: "/help",
        },
        {
          label: "Guides",
          items: [
            {
              label: "Installing Additional Packages",
              link: "/guides/installing-packages",
            },
            {
              label: "Adding Steps",
              link: "/guides/adding-steps",
            },
            {
              label: "Developing Locally",
              link: "/guides/developing-locally",
            },
            {
              label: "Running SeaPack in Production",
              link: "/guides/running-seapack-in-production",
            },
          ],
        },
        {
          label: "Configuration",
          items: [
            { label: "Configuration File", link: "/config/file" },
            {
              label: "Environment Variables",
              link: "/config/environment-variables",
            },
            { label: "Mise", link: "/config/mise" },
            { label: "Excluding Files", link: "/config/excluding-files" },
          ],
        },
        {
          label: "Languages",
          items: [
            { label: "Node", link: "/languages/node" },
            { label: "Python", link: "/languages/python" },
            { label: "Go", link: "/languages/golang" },
            { label: "Java", link: "/languages/java" },
            { label: "Deno", link: "/languages/deno" },
          ],
        },
        {
          label: "Deploying",
          items: [
            { label: "Sealos", link: "/deploying/sealos" },
            { label: "GitHub Actions", link: "/deploying/github-actions" },
          ],
        },
        {
          label: "Reference",
          items: [
            { label: "CLI Commands", link: "/reference/cli" },
            { label: "BuildKit Frontend", link: "/reference/frontend" },
          ],
        },
        {
          label: "Architecture",
          items: [
            { label: "High Level Overview", link: "/architecture/overview" },
            {
              label: "Package Resolution",
              link: "/architecture/package-resolution",
            },
            {
              label: "Secrets and Variables",
              link: "/architecture/secrets",
            },
            { label: "BuildKit Generation", link: "/architecture/buildkit" },
            { label: "Caching", link: "/architecture/caching" },
            { label: "User Config", link: "/architecture/user-config" },
          ],
        },
        {
          label: "Contributing",
          link: "/contributing",
        },
      ],
    }),
  ],

  vite: {
    plugins: [tailwindcss()],
  },
});
