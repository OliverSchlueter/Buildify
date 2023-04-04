<script lang="ts">
export default {
  props: [
    "baseLink",
    "id",
    "time",
    "hash",
    "message",
    "downloadLink"
  ],
  computed: {
    fullDownloadLink(){
        return this.baseLink + this.downloadLink;
    },
    shortHash() {
        return this.hash.substring(0, 7)
    },
    timeFormatted(){
        const t = new Date(this.time);
        return t.getDay() + "." + t.getMonth() + "." + t.getFullYear();
    }
  }
}
</script>

<template>
    <div class="build">
        <div class="descriptions">
            <p class="id"><a :href="baseLink + downloadLink">Build #{{ id }}</a></p>
            <div class="description">
                <p>
                    <span class="hash">{{ shortHash }}</span>
                    <span class="message">{{ message }}</span>
                </p>
            </div>
            <div class="description right">
                <span class="time">{{ timeFormatted }}</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
    .build{
        display: flex;
    }
    .id{
        background-color: var(--primary-color-dark);
        padding: 5px 15px;
        border-radius: 15px;
        margin-right: 10px;
        font-weight: 700;
    }

    .id a{
        color: inherit;
        text-decoration: none;
    }

    .descriptions{
        display: flex;
        width: 100%;
        /* justify-content: space-between; */
    }

    .description{
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-content: center;
    }
    
    .descriptions .right{
        margin-left: auto;
    }
    .hash{
        color: var(--primary-color);
        margin-right: 10px;
        font-family: 'Courier New', Courier, monospace;
    }

    .time{
        margin-left: auto;
    }
</style>