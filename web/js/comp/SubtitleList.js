import React, { Component } from 'react';

import subtitleStore from '../stores/subtitle_store'

import FileTags from './FileTags'
import ScoreTag from './ScoreTag'

class SubtitleList extends Component {
  constructor() {
    super()
    this.state = subtitleStore.getState()
    this.updateSubtitles = this.updateSubtitles.bind(this)
  }

  updateSubtitles() {
    console.log("Update subtitles", subtitleStore.getSubtitles())
    this.setState(subtitleStore.getState())
  }

  componentWillMount() {
    subtitleStore.on("change", this.updateSubtitles)
  }

  componentWillUnmount() {
    subtitleStore.removeListener("change", this.updateSubtitles)
  }

  render() {
    if (!this.state.subtitles || this.state.subtitles.length === 0) {
      return (
        <h2 className="center">Select a file to display subtitles</h2>
      )
    }

    let subs = this.state.subtitles

    subs = subs.filter((s) => s.language === this.state.lang)

    subs = subs.map((s) => {
      return (
        <li className="clearfix flex center">
          <ScoreTag score={s.score} />
          <div className="col ellipsis">{s.media.name}</div>
          <div className="">
            <FileTags media={s}/>
          </div>
          <div className="right">
            <button className="small">Download</button>
          </div>
        </li>
      )
    })

    return (
      <ul className="subtitle-list">{subs}</ul>
    )
  }
}

export default SubtitleList
