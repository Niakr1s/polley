import React from "react";
import ToMainPage from "../Links/ToMainPage";
import styles from './FourThousandFour.module.css';

export class FourThousandFour extends React.Component {
    render() {
        return (
            <div>
                <ToMainPage></ToMainPage>
                <h2 className={styles.notFound}>404 not found</h2>
            </div>
        )
    }
}