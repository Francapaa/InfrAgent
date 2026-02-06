import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

const PUBLIC_ROUTES = ['/login', '/auth'];
const PROTECTED_ROUTES = ['/dashboard', '/onboarding'];

function isPublicRoute(pathname: string): boolean {
  return PUBLIC_ROUTES.some(route => pathname.startsWith(route));
}

function isProtectedRoute(pathname: string): boolean {
  return PROTECTED_ROUTES.some(route => pathname.startsWith(route));
}

export function middleware(request: NextRequest): NextResponse {
  const { pathname } = request.nextUrl;
  const token = request.cookies.get('auth_token')?.value;

  // Permitir acceso a archivos estáticos y API routes
  if (
    pathname.startsWith('/_next') ||
    pathname.startsWith('/api') ||
    pathname.includes('.')
  ) {
    return NextResponse.next();
  }

  // Si hay token y está en /login, redirigir a onboarding
  if (token && pathname === '/login') {
    return NextResponse.redirect(new URL('/onboarding', request.url));
  }

  // si onboarding esta terminado y token está ==> DASHBOARD

  // Si no hay token y es ruta protegida, redirigir a login
  if (!token && isProtectedRoute(pathname)) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
