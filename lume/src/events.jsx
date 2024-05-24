import EventCard from "./_components/EventCard.jsx";

export const title = "Events";
export const layout = "base.njk";
export const date = "2012-12-31";
export const desc = "A list of the upcoming events that I plan to attend and what I'll do there.";

export default ({ events }) => (
    <>
        <h1 className="text-3xl mb-4">Events</h1>
        <p className="my-4">Where in the world is Xe Iaso?</p>
        {events.events === undefined ? (
            <p className="my-4">
                I don't have any events planned right now or my events API is down. Check back later or <a href="/contact">let me know</a> if you see this message in error!
            </p>
        ) : (
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                {events.events.map((event) => EventCard(event))}
            </div>
        )}

        <div className="my-4 prose dark:prose-invert max-w-full">
            <p>
                If you'd like me to speak at an event, please <a href="/contact">contact me</a>! I'm always looking for new opportunities to share my knowledge and experiences. I'm also available for interviews, podcasts, and other media appearances.
            </p>
            <p>Please note that all conferences and meetups I speak at require a publicly posted code of conduct.</p>
        </div>
    </>
);