export const title = "Donate";
export const layout = "base.njk";

const GitHubIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" className="icon icon-tabler icons-tabler-outline icon-tabler-brand-github"><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path d="M9 19c-4.3 1.4 -4.3 -2.5 -6 -3m12 5v-3.5c0 -1 .1 -1.4 -.5 -2c2.8 -.3 5.5 -1.4 5.5 -6a4.6 4.6 0 0 0 -1.3 -3.2a4.2 4.2 0 0 0 -.1 -3.2s-1.1 -.3 -3.5 1.3a12.3 12.3 0 0 0 -6.2 0c-2.4 -1.6 -3.5 -1.3 -3.5 -1.3a4.2 4.2 0 0 0 -.1 3.2a4.6 4.6 0 0 0 -1.3 3.2c0 4.6 2.7 5.7 5.5 6c-.6 .6 -.6 1.2 -.5 2v3.5" /></svg>
);

const PatreonIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" className="icon icon-tabler icons-tabler-outline icon-tabler-brand-patreon"><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path d="M20 8.408c-.003 -2.299 -1.746 -4.182 -3.79 -4.862c-2.54 -.844 -5.888 -.722 -8.312 .453c-2.939 1.425 -3.862 4.545 -3.896 7.656c-.028 2.559 .22 9.297 3.92 9.345c2.75 .036 3.159 -3.603 4.43 -5.356c.906 -1.247 2.071 -1.599 3.506 -1.963c2.465 -.627 4.146 -2.626 4.142 -5.273z" /></svg>
);

const CashIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" className="icon icon-tabler icons-tabler-outline icon-tabler-cash"><path stroke="none" d="M0 0h24v24H0z" fill="none" /><path d="M7 15h-3a1 1 0 0 1 -1 -1v-8a1 1 0 0 1 1 -1h12a1 1 0 0 1 1 1v3" /><path d="M7 9m0 1a1 1 0 0 1 1 -1h12a1 1 0 0 1 1 1v8a1 1 0 0 1 -1 1h-12a1 1 0 0 1 -1 -1z" /><path d="M12 14a2 2 0 1 0 4 0a2 2 0 0 0 -4 0" /></svg>
);

const DonationCard = ({ icon, backgroundImageUrl, name, link, accentColor }) => {
  return (
    // The main container is now only responsible for positioning, background, and animations.
    <div className="relative bg-gray-800 rounded-xl shadow-lg transform hover:scale-105 transition-transform duration-300 overflow-hidden motion-reduce:transform-none motion-reduce:transition-none">
      {/* Blurred Background Image */}
      <img
        src={backgroundImageUrl}
        alt="" // This is a decorative image, so the alt text is empty.
        className="absolute inset-0 w-full h-full object-cover filter blur-md opacity-20 z-0"
        onError={(e) => { e.target.style.display = 'none'; }}
      />

      {/* Card Content Wrapper */}
      {/* This new wrapper handles the content layout (padding, flexbox) and sits on top of the background. */}
      <div className="relative z-10 flex flex-col items-center text-center justify-center p-8 h-full">
        <h3 className="text-2xl font-bold text-white mb-4">{name}</h3>
        {/* This div sets the size and color for the SVG icon. */}
        <div className="w-16 h-16 text-white mb-4 flex-shrink-0">
          {icon}
        </div>
        <a
          href={link}
          target="_blank"
          rel="noopener noreferrer"
          className="text-white font-semibold py-3 px-8 rounded-lg transition-opacity duration-300 hover:opacity-80 motion-reduce:transition-none"
          style={{ backgroundColor: accentColor }}
        >
          Donate
        </a>
      </div>
    </div>
  );
};

export default () => {
  const donationOptions = [
    {
      id: 1,
      name: 'GitHub Sponsors',
      icon: <GitHubIcon />,
      backgroundImageUrl: 'https://files.xeiaso.net/hero/summer-walk.avif',
      link: 'https://github.com/sponsors/Xe',
      accentColor: '#635BFF',
    },
    {
      id: 2,
      name: "Liberapay",
      icon: <CashIcon />,
      backgroundImageUrl: 'https://files.xeiaso.net/hero/skrunkly-cherries.avif',
      link: 'https://liberapay.com/Xe/',
      accentColor: '#0070BA',
    },
    {
      id: 3,
      name: 'Patreon',
      icon: <PatreonIcon />,
      backgroundImageUrl: 'https://files.xeiaso.net/hero/airplane-side-sunset.avif',
      link: 'https://www.patreon.com/cadey',
      accentColor: '#6772E5',
    },
  ];

  return (
    <>
      <h1 className="text-3xl mb-4">Donate</h1>
      <p className="mb-4">
        Want to help me make things better? Your contribution helps me continue my work. Choose your preferred way to donate below.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 mb-4">
        {donationOptions.map((option) => (
          <DonationCard
            key={option.id}
            icon={option.icon}
            backgroundImageUrl={option.backgroundImageUrl}
            name={option.name}
            link={option.link}
            accentColor={option.accentColor}
          />
        ))}
      </div>

      <p className="mb-4">
        More options will be added in the future as facts and circumstances demand.
      </p>
    </>
  );
};
