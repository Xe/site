import EventCard from "./_components/EventCard.jsx";

export const title = "Events";
export const layout = "base.njk";
export const date = "2012-12-31";
export const desc = "A list of the upcoming events that I plan to attend and what I'll do there.";

export default ({ events }) => {
    if (events.events === undefined) {
        return (
            <>
                <h1 className="text-3xl mb-4">Events</h1>
                <p>
                    I don't have any events planned right now or my events API is down. Check back later!
                </p>
            </>
        );
    }

    return (
        <>
            <h1 className="text-3xl mb-4">Events</h1>

            <p className="my-4">Where in the world is Xe Iaso?</p>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                {events.events.map((event) => EventCard(event))}
            </div>
        </>
    );
}