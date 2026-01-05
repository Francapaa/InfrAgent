import { NextResponse } from "next/server";


export async function POST (req: Request) {
    const body = await req.json();

    const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/login`,{
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
        credentials: "include",
    })

    const data = await response.json(); 

    return NextResponse.json(data,{
        status: response.status,
        headers: {
            "set-cookie": response.headers.get("set-cookie") || "",
        }
    });
}
