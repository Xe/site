import BrandGithub from "https://deno.land/x/tabler_icons_tsx@0.0.3/tsx/brand-github.tsx";
import LemonIcon from "https://deno.land/x/tabler_icons_tsx@0.0.3/tsx/lemon-2.tsx";
import CheckIcon from "https://deno.land/x/tabler_icons_tsx@0.0.5/tsx/check.tsx";
import IconChevronRight from "https://deno.land/x/tabler_icons_tsx@0.0.3/tsx/chevron-right.tsx";
import IconAlarm from "https://deno.land/x/tabler_icons_tsx@0.0.3/tsx/alarm.tsx";
import IconZZZ from "https://deno.land/x/tabler_icons_tsx@0.0.5/tsx/zzz.tsx";
import IconArmchair from "https://deno.land/x/tabler_icons_tsx@0.0.3/tsx/armchair.tsx";

export const layout = "bare.njk";
export const title = "Alvis - Never get paged again";
export const year = "2023";
export const series = "techaro";

export const Header = ({ active }) => {
  const menus = [
    { name: "Home", href: "/" },
    { name: "Blog", href: "/blog/" },
    { name: "Contact", href: "/contact/" },
    { name: "Resume", href: "/resume/" },
    { name: "Talks", href: "/talks/" },
    { name: "VODs", href: "/vods/" },
  ];

  return (
    <div className="bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-slate-100 mb-4 w-full max-w-screen-lg py-6 px-8 flex flex-col md:flex-row gap-4">
      <div className="flex items-center flex-1">
        <LemonIcon className="inline-block" aria-hidden="true" />
        <div className="text-2xl ml-1 font-bold">
          Alvis
        </div>
      </div>
      <ul className="flex items-center gap-6">
        {menus.map((menu) => (
          <li>
            <a
              href={menu.href}
              className={menu.href === active ? " font-bold border-b-2" : ""}
            >
              {menu.name}
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
};

export const Footer = () => {
  const menus = [
    {
      title: "Documentation",
      children: [
        { name: "Getting Started", href: "/blog/alvis/" },
        { name: "Guide", href: "/blog/alvis/" },
        { name: "API", href: "/blog/alvis/" },
        { name: "Showcase", href: "/blog/alvis/" },
        { name: "Pricing", href: "/blog/alvis/" },
      ],
    },
    {
      title: "Community",
      children: [
        { name: "Forum", href: "/blog/alvis/" },
        { name: "Discord", href: "/blog/alvis/" },
      ],
    },
  ];

  return (
    <div class="bg-slate-200 dark:bg-slate-800 text-slate-900 dark:text-slate-100 mt-4 flex flex-col md:flex-row w-full max-w-screen-lg gap-8 md:gap-16 px-8 py-8 text-sm">
      <div class="flex-1">
        <div class="flex items-center gap-1">
          <LemonIcon class="inline-block" aria-hidden="true" />
          <div class="font-bold text-2xl">
            Alvis
          </div>
        </div>
        <div class="text-gray-500">
          Give ChatGPT root in prod. What could possibly go wrong?
        </div>
      </div>

      {menus.map((item) => (
        <div class="mb-4" key={item.title}>
          <div class="font-bold">{item.title}</div>
          <ul class="mt-2">
            {item.children.map((child) => (
              <li class="mt-2" key={child.name}>
                <a href={child.href}>
                  {child.name}
                </a>
              </li>
            ))}
          </ul>
        </div>
      ))}

      <div class="text-gray-500 space-y-2">
        <div class="text-xs">
          Copyright 2023 Xeserv.<br />
          All rights reserved.
        </div>

        <a
          href="https://github.com/Xe"
          class="inline-block hover:text-black"
          aria-label="GitHub"
        >
          <BrandGithub aria-hidden="true" />
        </a>
      </div>
    </div>
  );
};

export const Hero = () => {
  return (
    <div
      class="w-full my-8 py-8 flex px-8 py-10 min-h-[24em] justify-center items-center flex-col gap-8 bg-cover bg-center bg-no-repeat bg-slate-50 rounded-xl text-white"
      style="background-image:linear-gradient(rgba(0, 0, 40, 0.8),rgba(0, 0, 40, 0.8)), url('https://cdn.xeiaso.net/file/christine-static/hero/basketball.jpg');"
    >
      <div class="space-y-4 text-center">
        <h1 class="text-3xl inline-block font-bold">Alvis</h1>
        <p class="text-xl max-w-lg text-blue-100">
          Don't worry about being paged, Alvis will take care of it.
        </p>
      </div>

      <div class="flex flex-col md:flex-row items-center">
        <a
          href="/blog/alvis/"
          class="block mt-4 text-blue-600 visited:text-blue-600 cursor-pointer inline-flex items-center group text-blue-800 bg-white px-8 py-2 rounded-md hover:bg-blue-50 font-bold"
        >
          Get started{" "}
        </a>
        <a
          href="/blog/alvis/"
          class="block mt-4 transition-colors text-blue-400 visited:text-blue-400 cursor-pointer inline-flex items-center group px-4 py-2 hover:text-blue-100"
        >
          Learn more{" "}
          <IconChevronRight
            class="inline-block w-5 h-5 transition group-hover:translate-x-0.5"
            aria-hidden="true"
          />
        </a>
      </div>
    </div>
  );
};

export const Features = () => {
  const featureItems = [
    {
      icon: IconAlarm,
      description:
        "Automatically respond to production incidents with our industry-leading AI. Alvis will diagnose the problem and fix it.",
    },
    {
      icon: IconZZZ,
      description:
        "Leave yourself free to sleep, write, work on projects, relax, or spend time with your family. Alvis will take care of it.",
    },
    {
      icon: IconArmchair,
      description:
        "Get a detailed report for every incident Alvis handles. This helps you understand what went wrong and how to fix it.",
      link: "/landing/alvis/P0001/",
    },
  ];

  return (
    <div class="flex my-8 flex-col md:flex-row gap-8 p-8">
      {featureItems.map((item) => {
        return (
          <div class="flex-1 space-y-2">
            <div class="bg-blue-600 inline-block p-3 rounded-xl text-white">
              <item.icon class="w-10 h-10" aria-hidden="true" />
            </div>
            <p class="text-xl">
              {item.description}
            </p>

            {item.link &&
              (
                <a class="block" href={item.link}>
                  <p class="text-blue-600 cursor-pointer hover:underline inline-flex items-center group">
                    Read More{" "}
                    <IconChevronRight
                      class="inline-block w-5 h-5 transition group-hover:translate-x-0.5"
                      aria-hidden="true"
                    />
                  </p>
                </a>
              )}
          </div>
        );
      })}
    </div>
  );
};

export const Testimonials = () => {
  const testimonials = [
    {
      name: "James Dornick",
      position: "Lead SRE @ WritePulse",
      color: "bg-[#9d789b]",
      avatar: "https://cdn.xeiaso.net/avatar/d69caff2a9fd74b7069be6ece1a06ac5",
      quote:
        "Honestly, I wasn't able to spend time with my kid until I set up Alvis. Now family time isn't turned into work time!",
    },
    {
      name: "Samus Rhodes",
      position: "SRE @ Worcation",
      color: "bg-[#5b5faf]",
      avatar: "https://cdn.xeiaso.net/avatar/9108f6d8df326e20a4c1983a910cd952",
      quote:
        "I was able to get a promotion after setting up Alvis. I can't believe I didn't do it sooner!",
    },
    {
      name: "Maria Robotnik",
      position: "SRE @ G.U.N.",
      color: "bg-[#d61a9b]",
      avatar: "https://cdn.xeiaso.net/avatar/f9f35179d3331f30f92793b724c26cf9",
      quote:
        "Every time I get paged, I just tell Alvis to take care of it. I don't even have to look at my phone anymore!",
    },
  ];

  return (
    <section class="py-8 text-neutral-700 dark:text-neutral-300">
      <div class="mx-auto text-center md:max-w-xl lg:max-w-3xl">
        <h3 class="mb-6 text-3xl font-bold">Testimonials</h3>
        <p class="mb-6 pb-2 md:mb-12 md:pb-0">
          Don't believe us? Trust the voices of your peers. They obviously know
          what they're talking about, don't they?
        </p>
      </div>

      <div class="grid gap-6 text-center md:grid-cols-3">
        {testimonials.map((testimonial) => (
          <div>
            <div className="block rounded-lg bg-white shadow-lg dark:bg-neutral-700 dark:shadow-black/30 h-full">
              <div
                className={`h-28 overflow-hidden rounded-t-lg ${testimonial.color}`}
              >
              </div>
              <div className="mx-auto mt-12 w-24 overflow-hidden rounded-full border-2 border-white bg-white dark:border-neutral-800 dark:bg-neutral-800">
                <img
                  src={testimonial.avatar}
                />
              </div>
              <div className="p-6">
                <h4 className="mb-2 text-2xl font-semibold">
                  {testimonial.name}
                </h4>
                <div className="mb-2">{testimonial.position}</div>
                <hr />
                <p className="mt-4">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="currentColor"
                    class="inline-block h-7 w-7 pr-2"
                    viewBox="0 0 24 24"
                  >
                    <path d="M13 14.725c0-5.141 3.892-10.519 10-11.725l.984 2.126c-2.215.835-4.163 3.742-4.38 5.746 2.491.392 4.396 2.547 4.396 5.149 0 3.182-2.584 4.979-5.199 4.979-3.015 0-5.801-2.305-5.801-6.275zm-13 0c0-5.141 3.892-10.519 10-11.725l.984 2.126c-2.215.835-4.163 3.742-4.38 5.746 2.491.392 4.396 2.547 4.396 5.149 0 3.182-2.584 4.979-5.199 4.979-3.015 0-5.801-2.305-5.801-6.275z" />
                  </svg>
                  {testimonial.quote}
                </p>
              </div>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
};

export const FAQs = () => {
  const faqs = [
    {
      question: "How does Alvis work?",
      answer:
        "Alvis takes your production playbooks and turns them into diagnostic and remediation scripts. It then runs them when you get paged. If it's able to fix the problem, it closes the incident. Otherwise it escalates to you.",
    },
    {
      question: "What AI does Alvis use?",
      answer:
        "Alvis uses the latest and greatest in large language model technology. It's all developed in-house and is not a wrapper to other services. This ensures your data is safe and secure with Alvis.",
    },
    {
      question:
        "What is the connection between Alvis and the Rhadamanthus organization?",
      answer:
        "Alvis is a completly separate entity from Rhadamanthus. We are not affiliated with them in any way. We wish them luck with Project Exodus and their conflicts with the Saviorites.",
    },
    {
      question:
        "What happens when I go over the number of incidents in my plan?",
      answer:
        "Alvis will automatically escalate your incidents to the oncall engineer. This ensures that Alvis is fast and ready for everyone!",
    },
  ];
  return (
    <div className="mx-auto max-w-7xl px-6 py-12 sm:pt-16 lg:px-8 lg:py-40">
      <div className="lg:grid lg:grid-cols-12 lg:gap-8">
        <div className="lg:col-span-5">
          <h2 className="text-2xl font-bold leading-10 tracking-tight text-stone-950 dark:text-stone-50">
            Frequently asked questions
          </h2>
          <p className="mt-4 text-base leading-7">
            Can’t find the answer you’re looking for? Reach out to our{" "}
            <a
              href="#"
              className="font-semibold text-indigo-600 hover:text-indigo-500 visited:text-indigo-600"
            >
              customer support
            </a>{" "}
            team.
          </p>
        </div>
        <div className="mt-10 lg:col-span-7 lg:mt-0">
          <dl className="space-y-10">
            {faqs.map((faq) => (
              <div key={faq.question}>
                <dt className="text-base font-semibold leading-7 text-stone-800 dark:text-stone-400">
                  {faq.question}
                </dt>
                <dd className="mt-2 text-base leading-7">{faq.answer}</dd>
              </div>
            ))}
          </dl>
        </div>
      </div>
    </div>
  );
};

export const Pricing = () => {
  const tiers = [
    {
      name: "Individual",
      id: "tier-individual",
      href: "/blog/alvis/",
      priceMonthly: "$10",
      description: "The essentials to provide your best work for clients.",
      features: [
        "50 incidents included per month",
        "48 hour response time",
        "Basic analytics",
        "Help center access",
        "Tailscale integration",
      ],
      mostPopular: false,
    },
    {
      name: "Startup",
      id: "tier-startup",
      href: "/blog/alvis/",
      priceMonthly: "$32",
      description: "A plan that scales with your rapidly growing business.",
      features: [
        "Custom incident response steps",
        "500 incidents included per month",
        "Same-day response time",
        "Advanced analytics",
      ],
      mostPopular: true,
    },
    {
      name: "Enterprise",
      id: "tier-enterprise",
      href: "/blog/alvis/",
      priceMonthly: "$500",
      description: "Dedicated support and infrastructure for your company.",
      features: [
        "Artificial General Intelligence",
        "Unlimited incidents",
        "SOCII/HIPPA/PCI compliance",
        "Custom analytics",
        "Custom SLAs",
      ],
      mostPopular: false,
    },
  ];

  function classNames(...classes) {
    return classes.filter(Boolean).join(" ");
  }

  return (
    <div className="py-8 sm:py-12">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-4xl text-center">
          <h2 className="text-base font-semibold leading-7 text-indigo-600">
            Pricing
          </h2>
          <p className="mt-2 text-4xl font-bold tracking-tight text-gray-900 dark:text-gray-100 sm:text-5xl">
            Pricing plans for teams of&nbsp;all&nbsp;sizes
          </p>
        </div>
        <div className="isolate mx-auto mt-8 grid max-w-md grid-cols-1 gap-y-8 sm:mt-10 lg:mx-0 lg:max-w-none lg:grid-cols-3">
          {tiers.map((tier, tierIdx) => (
            <div
              key={tier.id}
              className={classNames(
                tier.mostPopular ? "lg:z-10 lg:rounded-b-none" : "lg:mt-8",
                tierIdx === 0 ? "lg:rounded-r-none" : "",
                tierIdx === tiers.length - 1 ? "lg:rounded-l-none" : "",
                "flex flex-col justify-between rounded-3xl bg-stone-100 dark:bg-stone-800 p-8 ring-1 ring-gray-200 xl:p-10",
              )}
            >
              <div>
                <div className="flex items-center justify-between gap-x-4">
                  <h3
                    id={tier.id}
                    className={classNames(
                      tier.mostPopular
                        ? "text-indigo-600 dark:text-indigo-200"
                        : "text-gray-900 dark:text-gray-100",
                      "text-lg font-semibold leading-8",
                    )}
                  >
                    {tier.name}
                  </h3>
                  {tier.mostPopular
                    ? (
                      <p className="rounded-full bg-indigo-600/10 px-2.5 py-1 text-xs font-semibold leading-5 text-indigo-600">
                        Most popular
                      </p>
                    )
                    : null}
                </div>
                <p className="mt-4 text-sm leading-6 text-gray-600 dark:text-gray-300">
                  {tier.description}
                </p>
                <p className="mt-6 flex items-baseline gap-x-1">
                  <span className="text-4xl font-bold tracking-tight text-gray-900 dark:text-gray-100">
                    {tier.priceMonthly}
                  </span>
                  <span className="text-sm font-semibold leading-6 text-gray-600 dark:text-gray-300">
                    /month
                  </span>
                </p>
                <ul
                  role="list"
                  className="mt-8 space-y-3 text-sm leading-6 text-gray-600 dark:text-gray-300"
                >
                  {tier.features.map((feature) => (
                    <li key={feature} className="flex gap-x-3">
                      <CheckIcon
                        className="h-6 w-5 flex-none stroke-indigo-600"
                        aria-hidden="true"
                      />
                      {feature}
                    </li>
                  ))}
                </ul>
              </div>
              <a
                href={tier.href}
                aria-describedby={tier.id}
                className={classNames(
                  tier.mostPopular
                    ? "bg-indigo-600 dark:bg-indigo-300 text-white dark:text-black shadow-sm hover:text-white visited:text-black hover:bg-indigo-500"
                    : "text-indigo-600 dark:text-indigo-600 ring-1 ring-inset ring-indigo-200 hover:ring-indigo-300",
                  "mt-8 block rounded-md py-2 px-3 text-center text-sm font-semibold leading-6 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600",
                )}
              >
                Buy plan
              </a>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export const CTA = () => {
  return (
    <div className="mx-auto max-w-7xl px-6 pb-8 sm:pb-10 lg:px-8">
      <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
        Boost your productivity.
        <br />
        Start using Alvis today.
      </h2>
      <div className="mt-10 flex items-center gap-x-6">
        <a
          href="/blog/alvis/"
          className="rounded-md visited:text-white bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
        >
          Get started
        </a>
        <a
          href="/blog/alvis/"
          className="text-sm font-semibold leading-6 text-gray-900"
        >
          Learn more <span aria-hidden="true">→</span>
        </a>
      </div>
    </div>
  );
};

export default function Alvis() {
  return (
    <body className="lg:max-w-5xl max-w-xl px-4 py-2 mx-auto dark:bg-stone-950 bg-stone-50 dark:text-stone-50 text-stone-950">
      <link
        rel="stylesheet"
        href="https://cdn.xeiaso.net/file/christine-static/static/font/inter/inter.css"
      />
      <div className="font-['Inter']">
        <Header active="/" />
        <Hero />
        <Features />
        <Testimonials />
        <Pricing />
        <FAQs />
        <CTA />
        <Footer />
      </div>
    </body>
  );
}
