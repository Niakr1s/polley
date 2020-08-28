import React from 'react'
import { IPoll } from '../../models/poll'
import styles from './Poll.module.css'
import { Formik, Form, Field } from 'formik'

interface IState {
    selected: number,
}

interface IProps {
    poll: IPoll,
    withVote: boolean, // indicates should vote selector and button be showed
    submitSelected: (selected: boolean[]) => void,
}

interface IFormValues {
    selected: boolean[],
}

interface IFormErrors {
    selected?: string,
}

class Poll extends React.Component<IProps, IState> {
    constructor(props: Readonly<IProps>) {
        super(props)
        this.state = {
            selected: 0,
        }
    }

    handleFormChange = (selected: boolean[]) => {
        let selectedCount = count(selected, true)
        if (selectedCount !== this.state.selected) {
            this.setState({ selected: selectedCount })
        }
    }

    formValidate = (values: IFormValues): IFormErrors => {
        let res: IFormErrors = {}
        if (choicesRemained(this.props.poll, count(values.selected, true)) < 0) res.selected = `only ${this.props.poll.allowMultiple} choices is allowed`
        return res
    }

    onFormSubmit = ({ selected }: IFormValues) => {
        this.props.submitSelected(selected)
    }

    render = () => {
        let choicesRem = choicesRemained(this.props.poll, this.state.selected)
        const pollExpired = isExpired(this.props.poll.expires)
        const voteAllowed = !pollExpired && this.props.poll.voteAllowed && this.props.withVote
        return (
            <div className={styles.pollMain}>
                <h2>{this.props.poll.name}</h2>

                <Formik initialValues={{ selected: Array(this.props.poll.choices.length).fill(false) } as IFormValues} onSubmit={this.onFormSubmit}
                    validate={this.formValidate}
                >
                    {({ values }) => {
                        this.handleFormChange(values.selected)
                        return (
                            <Form className={styles.pollContents}>
                                <table className={styles.table}>
                                    <tbody>
                                        {this.props.poll.choices.map((choice, idx) => (
                                            <tr key={choice.text}>
                                                <td>
                                                    <div className={styles.inline}>{choice.text}</div>
                                                </td>
                                                <td className={styles.progressBarContainer}>
                                                    <div className={styles.progressBar}
                                                        style={{ width: `${getPercentOfNthChoice(this.props.poll, idx)}%` }}
                                                    >{choice.votes}</div>
                                                    
                                                </td>
                                                <td>
                                                    {voteAllowed && <Field type="checkbox" name={`selected.${idx}`}></Field>}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                                {voteAllowed && <div className={styles.choicesRemained + " " + (choicesRem < 0 ? styles.red : "")}>Choices remained: {choicesRem}</div>}
                                {voteAllowed && <div>
                                    <input type="submit"></input>
                                </div>}
                            </Form>
                        )
                    }}
                </Formik>
                {pollExpired
                    ? <div className={styles.last}>Poll is expired.</div>
                    : <div className={styles.last}>{secondsRemained(this.props.poll.expires)} seconds remained</div>
                }
            </div>
        )
    }
}

const choicesRemained = (poll: IPoll, selected?: number): number => {
    return (poll.allowMultiple || 0) - (selected || 0)
}

const isExpired = (date: Date): boolean => {
    return date.valueOf() < Date.now()
}

const secondsRemained = (date: Date): number => {
    return Math.round((date.valueOf() - Date.now()) / 1000)
}

function count<T>(arr: T[], etalon: T): number {
    return arr.reduce((acc, v) => {
        if (v === etalon) acc++;
        return acc
    }, 0)
}

const getPercentOfNthChoice = (poll: IPoll, idx: number): number => {
    const sum: number = poll.choices.reduce((acc, choice) => acc + choice.votes, 0)
    const res: number = Math.round(poll.choices[idx].votes / sum * 100)
    return res
}

export default Poll