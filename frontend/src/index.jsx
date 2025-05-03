/* @refresh reload */
import { render } from 'solid-js/web';
import { Router, Route } from "@solidjs/router";

import { UserProvider } from './contexts/UserContext';

import './index.css';

import RouteGuard from './components/RouteGuard';

import Home from './home/Home';
import Contact from './contact/Contact';
import Pro from './pro/Pro';
import Profile from './profile/Profile';
import Terms from './terms/Terms';
import Header from './header/Header';
import Footer from './footer/Footer';

const root = document.getElementById("root");

if (!root) {
    throw new Error("Root div not found");
}


const Layout = (props) => {
    return (
        <>
            <Header />
            <RouteGuard>
                {props.children}
            </RouteGuard>
            <Footer />
        </>
    );
};

render(() => (
    <UserProvider>
        <Router root={Layout}>
            <Route path="/" component={Home} />
            <Route path="/contact" component={Contact} />
            <Route path="/pro" component={Pro} />
            <Route path="/profile" component={Profile} />
            <Route path="/terms" component={Terms} />
        </Router>
    </UserProvider>
), root);
