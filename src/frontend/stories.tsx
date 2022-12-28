import { g, u, x } from "xeact";

type Story = {
    steps: Step[],
};

type Step = {
    file: string,
    title: string,
    text: string,
};

const init = async (elemID: string, story: string) => {
    let root = g(elemID);
    x(root);
    let resp = await fetch(u(`/static/${story}.json`));
    if (!resp.ok) {
        root.appendChild(<div>
            <big>Oopsie-whoopsie!</big>
            <p>Oopsie-whoopsie uwu we made a fucky-wucky! A wittle fucko boingo! The code monkeys at our headquarters are working reawy hard to fix this!</p>
            <code>
                Wanted response to be ok, but unexpectedly got: {resp.status}
            </code>
        </div>);
        return;
    }

    let data: Story = await resp.json();
    console.log(data);

    const steps = data.steps.map(s => {
        console.log(s);
        //<article class="story" style="--bg: url(https://picsum.photos/480/840);"></article>
        const style = `--bg: url(https://cdn.xeiaso.net/file/christine-static/stories/${story}/${s.file}.jpg)`;
        console.log(style);
        const result = <article class="story" style={style}>
            <br />
            <br />
            <center><div class="wordart superhero"><span class="text">{s.title}</span></div></center>
            <br />
            <br />
            <div class="story-text text box">{s.text}</div>
        </article>;
        console.log(result.style.cssText);

        result.style.cssText = result.style.cssText.replace("\\", "");
        console.log(result.style.cssText);

        return result;
    });

    const stories = <div class="stories">
        <section class="user">
            {steps}
        </section>
    </div>;

    root.appendChild(stories);

    const median = stories.offsetLeft + (stories.clientWidth / 2);
    const state = {
        current_story: stories?.firstElementChild?.lastElementChild,
    };

    const navigateStories = (direction: "next" | "prev") => {
        const story = state.current_story;
        const lastItemInUserStory = story?.parentNode?.firstElementChild;
        const firstItemInUserStory = story?.parentNode?.lastElementChild;
        const hasNextUserStory = story?.parentElement?.nextElementSibling;
        const hasPrevUserStory = story?.parentElement?.previousElementSibling;

        if (direction === "next") {
            if (lastItemInUserStory === story && !hasNextUserStory) return;
            else if (lastItemInUserStory === story && hasNextUserStory) {
                state.current_story =
                    story.parentElement.nextElementSibling.lastElementChild;
                story?.parentElement.nextElementSibling.scrollIntoView({
                    behavior: "smooth",
                });
            } else {
                story?.classList.add("seen");
                state.current_story = story?.previousElementSibling;
            }
        } else if (direction === "prev") {
            if (firstItemInUserStory === story && !hasPrevUserStory) return;
            else if (firstItemInUserStory === story && hasPrevUserStory) {
                state.current_story =
                    story.parentElement.previousElementSibling.firstElementChild;
                story.parentElement.previousElementSibling.scrollIntoView({
                    behavior: "smooth",
                });
            } else {
                story?.nextElementSibling?.classList.remove("seen");
                state.current_story = story?.nextElementSibling;
            }
        }

        console.log(state.current_story.innerText);
    };

    root.addEventListener("click", (e) => {
        if (!(e.target instanceof HTMLElement)) {
            return;
        }
        if (e.target?.nodeName !== "ARTICLE") {
            return;
        }

        navigateStories(e.clientX > median ? "next" : "prev");
    });

    document.addEventListener("keydown", ({ key }) => {
        if (key === "ArrowDown" || key === "ArrowUp") {
            navigateStories(key === "ArrowDown" ? "next" : "prev");
        }
    });
};

export { init };
