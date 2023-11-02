const port = 8080;

const handler = async (req: Request): Promise<Response> => {
    const body = (await req.text());
    console.log(body);
    
    return new Response()
};

console.log(`HTTP server running. Access it at: http://localhost:${port}/`);
Deno.serve({ port }, handler);