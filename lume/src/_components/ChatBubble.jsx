export default function ChatBubble({
    reply = false,
    bg = "blue-dark",
    fg = "slate-50",
    children,
}) {
    return (
        <div className={`mx-auto 3xl:max-w-2xl ${reply ? "" : "space-y-4"}`}>
            <div className={`flex ${reply ? "justify-start" : "justify-end"}`}>
                <div className={`flex w-11/12 ${reply ? "" : "flex-row-reverse"}`}>
                    <div
                        className={`relative max-w-xl rounded-xl ${reply ? "rounded-tl-none" : "rounded-tr-none"
                            } bg-${bg} px-4 py-2`}
                    >
                        <span className={`font-medium text-${fg} font-['Inter']`}>{children}</span>
                    </div>
                </div>
            </div>
        </div>
    );
}