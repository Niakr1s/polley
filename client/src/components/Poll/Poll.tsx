import React from 'react'
import { IPoll, getPercentOfNthChoice, choicesRemained } from '../../models/poll'
import styles from './Poll.module.css'
import { Formik, Form, Field } from 'formik'
import { Link } from 'react-router-dom'
import { secondsRemained, isExpired } from '../../util/date'
import { count } from '../../util/arr'

interface IState {
    selected: number,
    secondsRemained?: number,
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

    secondsTimer?: ReturnType<typeof setInterval>

    startSeconds = (secondsRemained: number) => {
        this.setState({ secondsRemained })
        this.secondsTimer = setInterval(() => {
            if (this.state.secondsRemained == null) {
                this.stopSeconds();
            } else if (this.state.secondsRemained <= 0) {
                this.stopSeconds()
                this.setState({ secondsRemained: undefined })
            } else {
                this.setState({ secondsRemained: this.state.secondsRemained - 1 })
            }
        }, 1000)
    }

    stopSeconds = () => {
        this.secondsTimer && clearInterval(this.secondsTimer)
    }

    componentDidMount = () => {
        if (isExpired(this.props.poll.expires)) return
        this.startSeconds(secondsRemained(this.props.poll.expires))
    }

    componentWillUnmount = () => {
        this.stopSeconds()
    }

    render = () => {
        let choicesRem = choicesRemained(this.props.poll, this.state.selected)
        const pollExpired = isExpired(this.props.poll.expires)
        const voteAllowed = !pollExpired && this.props.poll.voteAllowed && this.props.withVote
        return (
            <div className="block">
                <h2>
                    <Link to={`/poll/${this.props.poll.uuid}`}>
                        {this.props.poll.name}
                    </Link>
                </h2>

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
                                                    ></div>
                                                </td>
                                                <td>{choice.votes}</td>
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
                {this.state.secondsRemained
                    ? <div className={styles.last}>{secondsRemained(this.props.poll.expires)} seconds remained</div>
                    : <div className={styles.last}>Poll is expired.</div>
                }
            </div>
        )
    }
}

export default Poll