const timestampPBtoDate = ({ seconds }) => {
    return new Date(seconds * 1000);
};

const formatDate = (date) => {
    return date.toLocaleDateString('en-US', {
        month: 'long',
        day: 'numeric',
        year: 'numeric',
    });
};

// takes within.website.x.mi.Event
export default ({ name, url, start_date, end_date, location, description }) => {
    const startDate = formatDate(timestampPBtoDate(start_date));
    const endDate = formatDate(timestampPBtoDate(end_date));
    return (
        <div className="rounded-lg p-4 bg-bg-1 dark:bg-bgDark-1">
            <h2 className="text-lg mb-2 text-fg-1 dark:text-fgDark-1">
                <a href={url} target="_blank" rel="noopener noreferrer" className="text-blue-dark dark:text-blueDark-light">
                    {name} <span role="img" aria-label="link">🔗</span>
                </a>
            </h2>
            <div className="card-content text-fg-1 dark:text-fgDark-1">
                <p>
                    {location} - {startDate} {start_date.seconds !== end_date.seconds ? `thru ${endDate}` : ""}
                </p>
                <p className="prose dark:prose-invert">
                    {description}
                </p>
            </div>
        </div>
    );
};
