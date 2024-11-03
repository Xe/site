export default function TecharoDisclaimer({ children }) {
    return (
        <>
            <link
                rel="stylesheet"
                href="https://cdn.xeiaso.net/file/christine-static/static/font/inter/inter.css"
            />
            <div className="font-['Inter'] text-lg mx-auto mt-4 mb-2 rounded-lg bg-bg-2 p-4 dark:bg-bgDark-2 md:max-w-3xl font-extrabold xe-dont-newline">
                {children}
            </div>
        </>
    );
}