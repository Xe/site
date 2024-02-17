export default function ChatFrame({ children }) {
    return (
        <>
            <div className="not-prose w-full space-y-4 p-4">{children}</div>
        </>
    );
}